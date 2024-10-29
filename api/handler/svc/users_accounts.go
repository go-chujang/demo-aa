package svc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/api/middleware/auth"
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/gofiber/fiber/v3"
)

type usersAccountsCreate struct {
	Owner common.Address `json:"owner"`
}

func UsersAccountsCreate(c fiber.Ctx) error {
	var req usersAccountsCreate
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}

	user, err := auth.UserAccount(c)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	if user.Owner != nil {
		return c.JSON(response.Body())
	}

	msg := model.CreateAccount{
		UserId: user.UserId,
		Owner:  req.Owner,
	}
	key, value, err := msg.KeyValue()
	if err != nil {
		return c.JSON(response.Err(err))
	}
	return c.JSON(response.Body().Err(kafka.Produce(msg.Topic(), key, value)))
}

func UsersAccountsState(c fiber.Ctx) error {
	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	owner, account, nonce, balances, err := aa.GetUserState(user)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	return c.JSON(response.Body().Data(map[string]interface{}{
		"owner":             *user.Owner,
		"ownerBalance":      hexutil.EncodeBig(owner),
		"account":           *user.Account,
		"accountBalance":    hexutil.EncodeBig(account),
		"accountBalances":   balances,
		"accountNonce":      hexutil.EncodeBig(nonce),
		"pending":           user.Pending,
		"lastFaucet":        user.LastFaucet,
		"syncedBlockNumber": user.SyncedBlockNumber,
	}))
}
