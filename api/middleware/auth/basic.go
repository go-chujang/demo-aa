package auth

import (
	"errors"

	"github.com/go-chujang/demo-aa/api/ctxutil"
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/internal/storedquery"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

const (
	owner_user     = "gochujang"
	owner_pass     = "demo"
	contextAccount = "_auth_basic_account"
)

func UserAccount(c fiber.Ctx) (*model.UserAccount, error) {
	if account, ok := c.Locals(contextAccount).(*model.UserAccount); !ok {
		return nil, errors.New("basic authentication failed or not provided")
	} else {
		return account, nil
	}
}

func UserAccountWithSync(c fiber.Ctx, requiredAccount ...bool) (*model.UserAccount, error) {
	account, err := UserAccount(c)
	if err != nil {
		return nil, err
	}
	switch account, err = aa.SyncUserState(account); {
	case err != nil:
		return nil, err
	case requiredAccount != nil && requiredAccount[0] && account.Account == nil:
		return nil, response.ErrByCode(response.ErrCodeAccountNotYet)
	default:
		return account, nil
	}
}

func Basic(db *mongox.Client) fiber.Handler {
	cfg := basicauth.ConfigDefault
	cfg.Unauthorized = func(c fiber.Ctx) error {
		c.Set(fiber.HeaderWWWAuthenticate, "basic realm="+cfg.Realm)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	cfg.Next = func(c fiber.Ctx) bool {
		return ctxutil.SkipAuth(c)
	}

	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		onlyOwner := ctxutil.OnlyOwner(c)
		if !onlyOwner && db == nil {
			return cfg.Unauthorized(c)
		}
		user, pass, ok := parseBasic(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return cfg.Unauthorized(c)
		}

		success := false
		switch {
		case onlyOwner:
			if user == owner_user && pass == owner_pass {
				success = true
			}
		default:
			account, err := storedquery.GetUser(db, user)
			if err != nil {
				return c.JSON(response.Err(err))
			}
			if user == account.UserId && pass == account.Password {
				c.Locals(contextAccount, account)
				success = true
			}
		}
		if success {
			return c.Next()
		}
		return cfg.Unauthorized(c)
	}
}
