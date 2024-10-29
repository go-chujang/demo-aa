package router

import (
	"errors"

	"github.com/go-chujang/demo-aa/api/middleware/auth"
	"github.com/go-chujang/demo-aa/api/middleware/logger"
	"github.com/go-chujang/demo-aa/api/middleware/routechart"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

type Router struct {
	cfg Config
	app *fiber.App
}

func New(config ...Config) (*Router, error) {
	cfg := configDefault(config...)

	app := fiber.New(cfg.fiberConfig())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	if *cfg.UseCors {
		app.Use(cors.New())
	}

	chartList := routechart.ChartWithDefault(cfg.RouteChart...)
	app.Use(routechart.New(chartList...))

	if *cfg.UseAuth {
		app.Use(auth.Basic(mongox.DB()))
	}

	app.Use(logger.NewResponseLogger(cfg.ResponseLogger))
	app.Use(logger.NewRequestLogger(cfg.RequestLogger))

	for _, chart := range chartList {
		group := app.Group(chart.Prefix)
		for _, row := range chart.Rows {
			handler := ternary.Cond(cfg.CtxTimeout > 0,
				timeout.New(row.Handler, cfg.CtxTimeout), row.Handler)

			switch row.Method {
			case fiber.MethodGet:
				group.Get(row.Path, handler)
			case fiber.MethodPost:
				group.Post(row.Path, handler)
			default:
				return nil, errors.New("unsupported method")
			}
		}
	}
	return &Router{cfg: cfg, app: app}, nil
}

func (r *Router) Start() error { return r.app.Listen(":" + r.cfg.Port) }
func (r *Router) Stop() error  { return r.app.Shutdown() }
