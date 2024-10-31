package routechart

import (
	"github.com/go-chujang/demo-aa/api/ctxutil"
	"github.com/go-chujang/demo-aa/api/handler/svc"
	"github.com/gofiber/fiber/v3"
)

type Chart struct {
	Prefix string
	Rows   []Row
}

type Row struct {
	Method  string
	Path    string
	Policy  *ctxutil.Policy
	Handler fiber.Handler
}

var (
	ServiceV1 = Chart{
		Prefix: "/svc/v1",
		Rows: []Row{
			{
				Method:  fiber.MethodPost,
				Path:    "/users/join",
				Policy:  ctxutil.NewPolicy().SkipAuth(),
				Handler: svc.UsersJoin,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/accounts/create",
				Handler: svc.UsersAccountsCreate,
			},
			{
				Method:  fiber.MethodGet,
				Path:    "/users/accounts/state",
				Handler: svc.UsersAccountsState,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/operations/execute",
				Handler: svc.UsersOperationsExecute,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/operations/tokens/transfer",
				Handler: svc.UsersOperationsTokensTransfer,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/operations/gambles/rock-paper-scissors",
				Handler: svc.UsersOperationsGamblesRockPaperScissors,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/operations/gambles/zero-to-nine",
				Handler: svc.UsersOperationsGamblesZeroToNine,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/users/operations/gambles/exchange",
				Handler: svc.UsersOperationsGamblesExchange,
			},
			{
				Method:  fiber.MethodPost,
				Path:    "/faucet",
				Handler: svc.Faucet,
			},
		},
	}
)
