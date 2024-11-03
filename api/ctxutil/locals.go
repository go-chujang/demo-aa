package ctxutil

import (
	"github.com/gofiber/fiber/v3"
)

func Locals[K, V any](c fiber.Ctx, key K) (V, bool) {
	value, ok := c.Locals(key).(V)
	return value, ok
}

func MustLocals[K, V any](c fiber.Ctx, key K) V {
	value, _ := c.Locals(key).(V)
	return value
}
