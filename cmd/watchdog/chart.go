package main

import (
	"errors"
	"sync"
	"time"

	"github.com/go-chujang/demo-aa/api/ctxutil"
	"github.com/go-chujang/demo-aa/api/middleware/routechart"
	"github.com/go-chujang/demo-aa/api/response"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/internal/storedquery"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"github.com/gofiber/fiber/v3"
)

func chart(db *mongox.Client, watchdog *aa.WatchDog, interval time.Duration, blockRange uint64) routechart.Chart {
	var (
		mu            sync.Mutex
		lastRequested = time.Now()
		limiter       = func() error {
			mu.Lock()
			defer mu.Unlock()
			now := time.Now()
			if now.Sub(lastRequested) < 3*time.Second {
				return errors.New("too frequently. wait for 3 seconds")
			}
			lastRequested = now
			return nil
		}
	)
	return routechart.Chart{
		Prefix: "/watchdog",
		Rows: []routechart.Row{
			{
				Method: fiber.MethodGet,
				Path:   "/status",
				Policy: ctxutil.NewPolicy().OnlyOwner(),
				Handler: func(c fiber.Ctx) error {
					blockNumber, _ := storedquery.BlockNumberInc(db, 0)
					return c.JSON(response.Body(map[string]interface{}{
						"envtag":        config.EnvTag(),
						"isStopped":     watchdog.IsStopped(),
						"interval":      interval,
						"blockRange":    blockRange,
						"blockNumber":   blockNumber,
						"lastRequested": lastRequested.Unix(),
					}))
				},
			},
			{
				Method: fiber.MethodPost,
				Path:   "/start",
				Policy: ctxutil.NewPolicy().OnlyOwner().SkipLog(),
				Handler: func(c fiber.Ctx) error {
					if err := limiter(); err != nil {
						return c.JSON(response.Err(err))
					}
					if !watchdog.IsStopped() {
						return c.JSON(response.Err(errors.New("already started")))
					}
					watchdog.Start()
					return c.JSON(response.Body("success"))
				},
			},
			{
				Method: fiber.MethodPost,
				Path:   "/stop",
				Policy: ctxutil.NewPolicy().OnlyOwner().SkipLog(),
				Handler: func(c fiber.Ctx) error {
					if err := limiter(); err != nil {
						return c.JSON(response.Err(err))
					}
					watchdog.Stop()
					return c.JSON(response.Body("success"))
				},
			},
		},
	}
}
