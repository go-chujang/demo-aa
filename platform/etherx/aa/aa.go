package aa

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chujang/demo-aa/common/utils/slice"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var fixedZeroSalt = big.NewInt(0)

type accountAbstract struct {
	db     *mongox.Client
	rpcUri string

	mngdWallet *model.MngdWallet
	caMap      map[common.Address]*etherx.Contractor

	entrypoint     *etherx.Contractor
	accountImpl    *etherx.Contractor
	accountFactory *etherx.Contractor
	tokenPaymaster *etherx.Contractor
	gamble         *etherx.Contractor
}

func newAA(db *mongox.Client, rpcUri string, readOnlyOps ...bool) (*accountAbstract, error) {
	var readOnly bool
	if readOnlyOps != nil && readOnlyOps[0] {
		readOnly = true
	}
	mngdWallet, entrypoint, accountImpl, accountFactory, tokenPaymaster, gamble, err := getMngdAccounts(db, rpcUri, readOnly)
	if err != nil {
		return nil, err
	}
	return &accountAbstract{
		db:         db,
		rpcUri:     rpcUri,
		mngdWallet: mngdWallet,
		caMap: map[common.Address]*etherx.Contractor{
			entrypoint.ContractAddress:     entrypoint,
			accountFactory.ContractAddress: accountFactory,
			tokenPaymaster.ContractAddress: tokenPaymaster,
			gamble.ContractAddress:         gamble,
		},
		entrypoint:     entrypoint,
		accountImpl:    accountImpl,
		accountFactory: accountFactory,
		tokenPaymaster: tokenPaymaster,
		gamble:         gamble,
	}, nil
}

func (a *accountAbstract) packCreateAccount(owner common.Address) ([]byte, error) {
	return a.accountFactory.Pack("createAccount", owner, fixedZeroSalt)
}

func (a *accountAbstract) packTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return a.tokenPaymaster.Pack("transfer", to, value)
}

func (a *accountAbstract) packBalanceOf(owner common.Address) ([]byte, error) {
	return a.tokenPaymaster.Pack("balanceOf", owner)
}

func (a *accountAbstract) unpackBalanceOf(hex string) (balance *big.Int, err error) {
	return balance, a.tokenPaymaster.Unpack(&balance, "balanceOf", hex)
}

func (a *accountAbstract) packCallData(target common.Address, data []byte) ([]byte, error) {
	return a.accountImpl.Pack("execute", target, big.NewInt(0), data)
}

func (a *accountAbstract) packGetNonce() ([]byte, error) {
	return a.accountImpl.Pack("getNonce")
}

func (a *accountAbstract) unpackGetNonce(hex string) (balance *big.Int, err error) {
	return balance, a.accountImpl.Unpack(&balance, "getNonce", hex)
}

func (a *accountAbstract) packHandleOps(beneficiary common.Address, ops ...PackedUserOperation) ([]byte, error) {
	return a.entrypoint.PackAny("handleOps", map[string]any{
		"ops":         ops,
		"beneficiary": beneficiary,
	})
}

var gambleTokenIds = []*big.Int{GambleTokenIdRockPaperScissors, GambleTokenIdZeroToNine}

func (a *accountAbstract) packBalanceOfBatch(account common.Address) ([]byte, error) {
	return a.gamble.Pack("balanceOfBatch", []common.Address{account, account}, gambleTokenIds)
}

func (a *accountAbstract) unpackBalanceOfBatch(hex string) (balances map[*big.Int]*big.Int, err error) {
	var unpacked []*big.Int
	if err = a.gamble.Unpack(&unpacked, "balanceOfBatch", hex); err != nil {
		return nil, err
	}
	return map[*big.Int]*big.Int{
		GambleTokenIdRockPaperScissors: unpacked[0],
		GambleTokenIdZeroToNine:        unpacked[1],
	}, nil
}

func (a *accountAbstract) packRockPaperScissors(choice string) ([]byte, error) {
	return a.gamble.Pack("rockPaperScissors", choice)
}

func (a *accountAbstract) packZeroToNine(guess *big.Int) ([]byte, error) {
	return a.gamble.Pack("zeroToNine", guess)
}

func (a *accountAbstract) packExchange(tokenId, amount *big.Int) ([]byte, error) {
	return a.gamble.Pack("exchange", tokenId, amount)
}

// ca must exist in aa.caMap
type BundlePayload struct {
	ca     common.Address
	packed []byte
}

func (a *accountAbstract) bundleExec(txr *etherx.Transactor, payloads []BundlePayload) (hashes []string, rpcErrs []error, err error) {
	nonce, gasPrice, err := txr.NonceAndGasPrice()
	if err != nil {
		return nil, nil, err
	}
	txns := slice.TypeCastWithIdx(payloads, func(i int, bp BundlePayload) *types.Transaction {
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce + uint64(i),
			Gas:      a.caMap[bp.ca].GetMaxGasLimit(),
			GasPrice: gasPrice,
			To:       &bp.ca,
			Data:     bp.packed,
		})
	})
	return txr.SendTransactionsWithSign(txns)
}

func (a *accountAbstract) bunleCommit(ctx context.Context,
	hashes []string, rpcErrs []error, rawDataHelpers []model.RawDataHelper) ([]error, error) {

	writeModels := make([]mongo.WriteModel, len(hashes))
	for i, hash := range hashes {
		hint, rawData := rawDataHelpers[i].RawData()
		updateSet := bson.M{
			"hint":    hint,
			"rawData": rawData,
		}
		if rpcErrs[i] != nil {
			updateSet["failedReason"] = rpcErrs[i].Error()
		}
		writeModels[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": hash}).
			SetUpdate(bson.M{"$set": updateSet}).
			SetUpsert(true)
	}
	query := mongox.NewQuery().SetColl(model.CollectionTxLogs).SetOrdered(false)
	return a.db.BulkWrite(writeModels, query, ctx)
}

func (a *accountAbstract) getDeposit() (value *big.Int, err error) {
	return value, a.tokenPaymaster.Call(&value, "getDeposit")
}

func (a *accountAbstract) deposit(depositor *etherx.Transactor, value *big.Int) (string, error) {
	packed, err := a.tokenPaymaster.Pack("deposit")
	if err != nil {
		return "", err
	}
	nonce, gasPrice, err := depositor.NonceAndGasPrice()
	if err != nil {
		return "", err
	}
	hashes, rpcErrs, err := depositor.SendTransactionsWithSign([]*types.Transaction{types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      50_000,
		GasPrice: gasPrice,
		To:       &a.tokenPaymaster.ContractAddress,
		Data:     packed,
		Value:    value,
	})})
	if err != nil {
		return "", err
	}
	if err = errors.Join(rpcErrs...); err != nil {
		return "", err
	}
	return hashes[0], nil
}

func (a *accountAbstract) balanceOf(owner common.Address) (balance *big.Int, err error) {
	return balance, a.tokenPaymaster.Call(&balance, "balanceOf", owner)
}

func (a *accountAbstract) balanceOfBatch(owner common.Address) (balances []*big.Int, err error) {
	return balances, a.tokenPaymaster.Call(&balances, "balanceOfBatch", owner)
}

func (a *accountAbstract) getAddress(owner common.Address) (addr common.Address, err error) {
	return addr, a.accountFactory.Call(&addr, "getAddress", owner, fixedZeroSalt)
}

func (a *accountAbstract) getNonce(accountAddress common.Address) (nonce *big.Int, err error) {
	pack, err := a.packGetNonce()
	if err != nil {
		return nil, err
	}
	res, err := rpcx.DoReq[string](a.rpcUri, a.accountImpl.CallReq(pack, 0, accountAddress))
	if err != nil {
		return nil, err
	}
	return a.unpackGetNonce(res)
}

func (a *accountAbstract) getUserOpHash(userOp PackedUserOperation) (hash [32]byte, err error) {
	return hash, a.entrypoint.CallAny(&hash, "getUserOpHash", map[string]any{"userOp": userOp})
}
