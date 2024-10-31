package router

import (
	"io"
	"os"
	"time"

	"github.com/go-chujang/demo-aa/api/middleware/routechart"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"github.com/go-chujang/demo-aa/internal/json"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	ForceSet *fiber.Config

	AppName        string
	Port           string
	CtxTimeout     time.Duration
	Concurrency    int
	ReleaseMode    bool
	RequestLogger  io.Writer
	ResponseLogger io.Writer
	RWTimeout      time.Duration

	// middlewares
	UseLiveness     *bool
	UseReadiness    *bool
	UseCors         *bool
	UseCacheControl *bool
	UseAuth         *bool

	RouteChart []routechart.Chart
}

var ConfigDefault = Config{
	ForceSet:        nil,
	AppName:         "demo",
	Port:            "5000",
	CtxTimeout:      time.Second * 3,
	Concurrency:     fiber.DefaultConcurrency,
	ReleaseMode:     false,
	RWTimeout:       time.Second * 2,
	RequestLogger:   os.Stdout,
	ResponseLogger:  os.Stdout,
	UseLiveness:     conv.ToPtr(true),
	UseReadiness:    conv.ToPtr(false),
	UseCors:         conv.ToPtr(false),
	UseCacheControl: conv.ToPtr(false),
	UseAuth:         conv.ToPtr(true),
}

func configDefault(config ...Config) Config {
	var cfg Config
	if config == nil {
		cfg = ConfigDefault
	} else {
		cfg = config[0]
	}

	if cfg.AppName == "" {
		cfg.AppName = ConfigDefault.AppName
	}
	if cfg.Port == "" {
		cfg.Port = ConfigDefault.Port
	}
	if cfg.CtxTimeout <= 0 {
		cfg.CtxTimeout = ConfigDefault.CtxTimeout
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = ConfigDefault.Concurrency
	}
	if cfg.RWTimeout <= 0 {
		cfg.RWTimeout = ConfigDefault.RWTimeout
	}
	if cfg.RequestLogger == nil {
		cfg.RequestLogger = ConfigDefault.RequestLogger
	}
	if cfg.ResponseLogger == nil {
		cfg.ResponseLogger = ConfigDefault.ResponseLogger
	}
	if cfg.UseLiveness == nil {
		cfg.UseLiveness = ConfigDefault.UseLiveness
	}
	if cfg.UseReadiness == nil {
		cfg.UseReadiness = ConfigDefault.UseReadiness
	}
	if cfg.UseCors == nil {
		cfg.UseCors = ConfigDefault.UseCors
	}
	if cfg.UseCacheControl == nil {
		cfg.UseCacheControl = ConfigDefault.UseCacheControl
	}
	if cfg.UseAuth == nil {
		cfg.UseAuth = ConfigDefault.UseAuth
	}
	return cfg
}

func (c Config) fiberConfig() fiber.Config {
	if c.ForceSet != nil {
		return *c.ForceSet
	}

	return fiber.Config{
		AppName:                  c.AppName,
		Concurrency:              c.Concurrency,
		ReadTimeout:              c.RWTimeout,
		WriteTimeout:             c.RWTimeout,
		CaseSensitive:            false,
		StrictRouting:            true,
		Immutable:                true,
		DisableKeepalive:         true,
		DisableHeaderNormalizing: false,
		ReduceMemoryUsage:        false,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	}
}
