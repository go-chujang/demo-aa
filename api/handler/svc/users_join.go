package svc

import (
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type usersJoin struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

func UsersJoin(c fiber.Ctx) error {
	var req usersJoin
	if err := c.Bind().Body(&req); err != nil {
		return c.JSON(response.Err(err))
	}
	user := model.UserAccount{
		UserId:   req.UserId,
		Password: req.Password,
	}
	switch err := mongox.DB().InsertOne(user, mongox.NewQuery().SetColl(user.Collection())); {
	case err != nil:
		return c.JSON(response.Err(err))
	case mongo.IsDuplicateKeyError(err):
		return c.JSON(response.Body().Code(response.ErrAlreadyExistUserId))
	default:
		return c.JSON(response.Body().Err(err))
	}
}
