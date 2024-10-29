package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

func setup_aa() {
	var (
		mngdAccounts          []model.ManagedAccount
		deployer              *etherx.Transactor
		deployerAddr          common.Address
		manager               common.Address
		keys, _               = ethutil.ParseMnemonic(mnemonic, uint32(accountCnt))
		premintAmount, _      = new(big.Int).SetString("10000000000000000000000000000", 10) // 10 billion
		depositAmount, _      = new(big.Int).SetString("10000000000000000000", 10)          // 10 eth
		toManagerAmount, _    = new(big.Int).SetString("100000000000000000000000000", 10)   // 100 million token
		toGambleRewardPool, _ = new(big.Int).SetString("10000000000000000000", 10)          // 10 token
	)
	for i, key := range keys {
		var (
			role    = model.RoleOperator
			addr, _ = ethutil.PvKey2Address(key)
		)
		switch i {
		case 0:
			role = model.RoleSupervisor
			deployer, _ = etherx.NewTransactor(rpcuri, key)
			deployerAddr = addr
		case 1:
			role = model.RoleManager
			manager = addr
		}
		mngdAccounts = append(mngdAccounts, model.ManagedAccount{
			Id:         fmt.Sprintf("manager%d", i),
			Kind:       model.KindEOA,
			Role:       &role,
			Address:    addr,
			PrivateKey: key,
		})
	}
	if deployer == nil {
		panic("something wrong")
	}

	// Contracts
	appendCA := func(deployFn func(bytecode []byte, args ...interface{}) (string, common.Address, error),
		name string, bytecode []byte, args ...interface{}) {
		hash, ca, err := deployFn(bytecode, args...)
		if err != nil {
			panic(err)
		}
		mngdAccounts = append(mngdAccounts, model.ManagedAccount{
			Id:             name,
			Kind:           model.KindContract,
			Address:        ca,
			DeployedTxHash: &hash,
			Deployer:       &deployer.Address,
		})
	}
	entrypoint, err := etherx.NewContractor(rpcuri, aa.EntrypointAbi, ethutil.ZeroAddress, deployer)
	if err != nil {
		panic(err)
	}
	appendCA(entrypoint.Deploy, aa.EntrypointName.String(), aa.EntrypointBytes)

	accountFactory, err := etherx.NewContractor(rpcuri, aa.AccountFactoryAbi, ethutil.ZeroAddress, deployer)
	if err != nil {
		panic(err)
	}
	appendCA(accountFactory.Deploy, aa.AccountFactoryName.String(), aa.AccountFactoryBytes,
		entrypoint.ContractAddress)

	paymaster, err := etherx.NewContractor(rpcuri, aa.PaymasterAbi, ethutil.ZeroAddress, deployer)
	if err != nil {
		panic(err)
	}
	appendCA(paymaster.Deploy, aa.PaymasterName.String(), aa.PaymasterBytes,
		accountFactory.ContractAddress, "DEMO", entrypoint.ContractAddress)

	gamble, err := etherx.NewContractor(rpcuri, aa.GambleAbi, ethutil.ZeroAddress, deployer)
	if err != nil {
		panic(err)
	}
	appendCA(gamble.Deploy, aa.GambleName.String(), aa.GambleBytes, paymaster.ContractAddress)

	// deposit
	packed, err := paymaster.Pack("deposit")
	if err != nil {
		panic(err)
	}
	depositGas := uint64(100_000)
	txns, err := deployer.ValueTxs(&depositGas, []*common.Address{&paymaster.ContractAddress}, []*big.Int{depositAmount}, packed)
	if err != nil {
		panic(err)
	}
	_, rpcErrs, err := deployer.SendTransactionsWithSign(txns)
	if err != nil {
		panic(err)
	}
	if err = errors.Join(rpcErrs...); err != nil {
		panic(err)
	}

	for i, doc := range mngdAccounts {
		if err = db.InsertOne(&doc, mongox.NewQuery().SetColl(doc.Collection())); err != nil {
			panic(fmt.Errorf("slice idx: %d, err: %s", i, err.Error()))
		}
	}
	if _, err := paymaster.ExecuteOne("mintTokens", deployerAddr, premintAmount); err != nil {
		panic("failed premint")
	}
	// token transfer to manager
	if _, err := paymaster.ExecuteOne("transfer", manager, toManagerAmount); err != nil {
		panic("failed transfer")
	}
	// token transfer to gamble-contract
	if _, err := paymaster.ExecuteOne("transfer", gamble.ContractAddress, toGambleRewardPool); err != nil {
		panic("failed transfer")
	}
}
