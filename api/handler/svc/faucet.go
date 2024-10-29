package svc

import (
	"time"

	"github.com/go-chujang/demo-aa/api/middleware/auth"
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	faucetAmount, _ = config.GetBig(config.FAUCET_AMOUNT)
	faucetLimit, _  = config.GetDurationSeconds(config.FAUCET_LIMIT)
)

// todo: update lastFaucet
func Faucet(c fiber.Ctx) error {
	user, err := auth.UserAccountWithSync(c, true)
	if err != nil {
		return c.JSON(response.Err(err))
	}

	now := time.Now().Unix()
	if user.LastFaucet+int64(faucetLimit) > now {
		return c.JSON(response.Body().Code(response.ErrFaucetLimit))
	}
	msg := model.Faucet{
		Receiver: *user.Account,
		Value:    faucetAmount,
	}
	key, value, err := msg.KeyValue()
	if err != nil {
		return c.JSON(response.Err(err))
	}
	if err = kafka.Produce(msg.Topic(), key, value); err != nil {
		return c.JSON(response.Err(err))
	}
	err = mongox.DB().UpdateOne(mongox.NewQuery().
		SetColl(model.CollectionUserAccounts).
		SetFilter(bson.M{"_id": user.UserId}).
		SetUpdate(bson.M{"$set": bson.M{"lastFaucet": now}}))
	return c.JSON(response.Body().Err(err))
}
