package aa

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

var (
	api     *accountAbstract
	apiOnce sync.Once
	caMap   map[contractName]common.Address
)

func apiInstanceMustNotNil() {
	if api == nil {
		panic("api insatance must not nil, UseApi() first")
	}
}

func UseApi(db *mongox.Client, rpcUri string) (err error) {
	apiOnce.Do(func() {
		api, err = newAA(db, rpcUri, true)
		caMap = map[contractName]common.Address{
			EntrypointName:     api.entrypoint.ContractAddress,
			AccountFactoryName: api.accountFactory.ContractAddress,
			PaymasterName:      api.tokenPaymaster.ContractAddress,
			GambleName:         api.gamble.ContractAddress,
		}
	})
	return err
}

var (
	gambleZeroBalances = map[*big.Int]*big.Int{
		GambleTokenIdRockPaperScissors: big.NewInt(0),
		GambleTokenIdZeroToNine:        big.NewInt(0),
	}
)

func GetUserState(user *model.UserAccount) (
	ownerBalance *big.Int,
	accountBalance *big.Int,
	accountNonce *big.Int,
	accountBalances map[*big.Int]*big.Int,
	err error) {

	if user.Owner == nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), gambleZeroBalances, nil
	}
	pack0, _ := api.packBalanceOf(*user.Owner)
	pack1, _ := api.packBalanceOf(*user.Account)
	pack2, _ := api.packGetNonce()
	pack3, _ := api.packBalanceOfBatch(*user.Account)

	batch := []rpcx.Request{
		api.tokenPaymaster.CallReq(pack0, 0),
		api.tokenPaymaster.CallReq(pack1, 1),
		api.accountImpl.CallReq(pack2, 2, *user.Account),
		api.gamble.CallReq(pack3, 3),
	}
	_, res, _, err := rpcx.EasyBatch(api.rpcUri, batch)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if len(res) != 4 {
		return nil, nil, nil, nil, errors.New("todo debug")
	}
	res0, err0 := api.unpackBalanceOf(res[0])
	res1, err1 := api.unpackBalanceOf(res[1])
	res2, err2 := api.unpackGetNonce(res[2])
	res3, err3 := api.unpackBalanceOfBatch(res[3])
	return res0, res1, res2, res3, errors.Join(err0, err1, err2, err3)
}

func BalanceOf(account common.Address) (*big.Int, error) {
	return api.balanceOf(account)
}

func BalanceOfBatch(account common.Address) (rockPaperScissors *big.Int, zeroToNine *big.Int, err error) {
	balances, err := api.balanceOfBatch(account)
	if err != nil {
		return nil, nil, err
	}
	return balances[0], balances[1], nil
}

func SyncUserState(before *model.UserAccount) (*model.UserAccount, error) {
	apiInstanceMustNotNil()

	if before.Account != nil && !before.Pending {
		return before, nil
	}
	if before.Account == nil {
		return syncCreateAccount(before)
	} else {
		return syncUserOperation(before)
	}
}

func userOperation(target contractName, user *model.UserAccount, packed []byte) (*UserOperationWithHash, error) {
	targetCA, exist := caMap[target]
	if !exist {
		return nil, errors.New("unsupported target")
	}

	account := *user.Account
	nonce, err := api.getNonce(account)
	if err != nil {
		return nil, err
	}
	callData, err := api.packCallData(targetCA, packed)
	if err != nil {
		return nil, err
	}

	paymaster := caMap[PaymasterName]
	op := userOperationDefault(paymaster).
		SetSender(account).
		SetNonce(nonce).
		SetCallData(callData)
	userOpHash, err := api.getUserOpHash(*op)
	if err != nil {
		return nil, err
	}
	return &UserOperationWithHash{
		UserOperation: op,
		Hash:          userOpHash[:],
	}, nil
}

func VerifySignature(user *model.UserAccount, userOperation PackedUserOperation) (bool, error) {
	owner := *user.Owner
	sig := make([]byte, len(userOperation.Signature))
	copy(sig, userOperation.Signature)

	copied := userOperation
	copied.SetSignature(nil)

	hash, err := api.getUserOpHash(copied)
	if err != nil {
		return false, err
	}
	return ethutil.VerifySignature(owner, hash[:], sig)
}

func UserOpTransferToken(user *model.UserAccount, to common.Address, value *big.Int) (*UserOperationWithHash, error) {
	packed, err := api.packTransfer(to, value)
	if err != nil {
		return nil, err
	}
	return userOperation(PaymasterName, user, packed)
}

func GambleRockPaperScissors(user *model.UserAccount, choice string) (*UserOperationWithHash, error) {
	packed, err := api.packRockPaperScissors(choice)
	if err != nil {
		return nil, err
	}
	return userOperation(GambleName, user, packed)
}

func GambleZeroToNince(user *model.UserAccount, guess *big.Int) (*UserOperationWithHash, error) {
	packed, err := api.packZeroToNine(guess)
	if err != nil {
		return nil, err
	}
	return userOperation(GambleName, user, packed)
}

func GambleExchange(user *model.UserAccount, tokenId, amount *big.Int) (*UserOperationWithHash, error) {
	packed, err := api.packExchange(tokenId, amount)
	if err != nil {
		return nil, err
	}
	return userOperation(GambleName, user, packed)
}
