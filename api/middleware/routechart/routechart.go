package routechart

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type config struct {
	Next      func(c fiber.Ctx) bool
	NotFound  fiber.Handler
	routeInfo *node
}

func New(charts ...Chart) fiber.Handler {
	cfg := config{
		Next: nil,
		NotFound: func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotFound)
		},
		routeInfo: newNode(),
	}
	for _, chart := range charts {
		prefix := chart.Prefix
		for _, v := range chart.Rows {
			path := prefix + v.Path
			segs := strings.Split(path, "/")[1:]
			cfg.routeInfo.insert(segs, v.Policy)
		}
	}
	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		var (
			reqPath = c.Path()
			segs    = strings.Split(reqPath, "/")[1:]
		)
		if policy, exist := cfg.routeInfo.search(segs); exist && policy != nil {
			policy.Stores(c)
		}
		return c.Next()
	}
}
