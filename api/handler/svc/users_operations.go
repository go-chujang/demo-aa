package svc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/api/middleware/auth"
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/gofiber/fiber/v3"
)

type usersOperationsExecute struct {
	Op *aa.PackedUserOperation `json:"op"`
}

func UsersOperationsExecute(c fiber.Ctx) error {
	var req usersOperationsExecute
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}

	// todo: gas calc
	balance, err := aa.BalanceOf(req.Op.Sender)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	if len(balance.Bits()) == 0 {
		return c.JSON(response.Body().Code(response.ErrTxnInsufficientToken))
	}

	switch user, err := auth.UserAccountWithSync(c, true); {
	case err != nil:
		return c.JSON(response.Err(err))
	case user.Pending:
		return c.JSON(response.Body().Code(response.ErrTxnPendingState))
	case req.Op.Nonce.Uint64() != 0 && req.Op.Nonce.Uint64() < user.LastUsedNonce:
		return c.JSON(response.Body().Code(response.ErrTxnNonceTooLow))
	default:
		if valid, err := aa.VerifySignature(user, *req.Op); err != nil {
			return c.JSON(response.Err(err))
		} else if !valid {
			return c.JSON(response.Body().Code(response.ErrTxnInvalidSig))
		}
	}
	msg := req.Op.Message()
	key, value, err := msg.KeyValue()
	if err != nil {
		return c.JSON(response.Err(err))
	}
	return c.JSON(response.Body().Err(kafka.Produce(msg.Topic(), key, value)))
}

type usersOperationsTokensTransfer struct {
	To    common.Address `json:"to"`
	Value string         `json:"value"`
}

func UsersOperationsTokensTransfer(c fiber.Ctx) error {
	var req usersOperationsTokensTransfer
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}
	if req.To.Cmp(ethutil.ZeroAddress) == 0 {
		return c.JSON(response.Body().Code(response.ErrInvalidAddress))
	}
	value, err := ethutil.HexOrString2X[*big.Int](req.Value)
	if err != nil {
		return c.JSON(response.Err(err))
	}

	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	userOp, err := aa.UserOpTransferToken(user, req.To, value)
	return c.JSON(response.Body(userOp).Err(err))
}

type usersOperationsGamblesRockPaperScissors struct {
	Choice string `json:"choice"`
}

func UsersOperationsGamblesRockPaperScissors(c fiber.Ctx) error {
	var req usersOperationsGamblesRockPaperScissors
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}
	switch req.Choice {
	case "rock":
	case "paper":
	case "scissors":
	default:
		return c.JSON(response.Body().Err(errors.New("invalid choice, must be rock-paper-scissors")))
	}
	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	userOp, err := aa.GambleRockPaperScissors(user, req.Choice)
	return c.JSON(response.Body(userOp).Err(err))
}

type usersOperationsGamblesZeroToNine struct {
	Guess uint64 `json:"guess"`
}

func UsersOperationsGamblesZeroToNine(c fiber.Ctx) error {
	var req usersOperationsGamblesZeroToNine
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}
	if req.Guess >= 10 {
		return c.JSON(response.Body().Err(errors.New("guess must be between 0 and 9")))
	}
	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}
	userOp, err := aa.GambleZeroToNince(user, new(big.Int).SetUint64(req.Guess))
	return c.JSON(response.Body(userOp).Err(err))
}

type usersOperationsGamblesExchange struct {
	TokenId string `json:"tokenId"`
}

func UsersOperationsGamblesExchange(c fiber.Ctx) error {
	var req usersOperationsGamblesExchange
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}
	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}

	tokenId, err := ethutil.HexOrString2X[*big.Int](req.TokenId)
	if err != nil {
		return c.JSON(response.Err(err))
	}

	var amount *big.Int
	switch rps, z2n, err := aa.BalanceOfBatch(*user.Account); {
	case err != nil:
		return c.JSON(response.Err(err))
	case tokenId.Cmp(aa.GambleTokenIdRockPaperScissors) == 0:
		if len(rps.Bits()) == 0 {
			return c.JSON(response.Body().Code(response.ErrTxnInsufficientToken))
		}
		amount = rps
	case tokenId.Cmp(aa.GambleTokenIdZeroToNine) == 0:
		if len(z2n.Bits()) == 0 {
			return c.JSON(response.Body().Code(response.ErrTxnInsufficientToken))
		}
		amount = z2n
	default:
		return c.JSON(response.Body().Err(fmt.Errorf("tokenId must be %d or %d", aa.GambleTokenIdRockPaperScissors.Uint64(), aa.GambleTokenIdZeroToNine.Uint64())))
	}

	userOp, err := aa.GambleExchange(user, tokenId, amount)
	return c.JSON(response.Body(userOp).Err(err))
}
