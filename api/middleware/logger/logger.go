package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/go-chujang/demo-aa/api/ctxutil"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const (
	timeFormat = time.RFC3339
	timeZone   = "UTC"
)

func NewRequestLogger(w io.Writer) fiber.Handler {
	return logger.New(logger.Config{
		Next: func(c fiber.Ctx) bool {
			return ctxutil.SkipRequestLog(c)
		},
		Output:     w,
		TimeFormat: timeFormat,
		TimeZone:   timeZone,
		Format: fmt.Sprintf("${%s} [%s] ${%s} ${%s} ${%s} ${%s} | ${%s} | ${%s%s} | ${%s}\n",
			logger.TagTime,
			"Request Log",
			logger.TagStatus,
			logger.TagMethod,
			logger.TagPath,
			logger.TagQueryStringParams,
			logger.TagBody,
			// ${locals:requestid}
			logger.TagLocals,
			"${respHeader:X-Request-ID}",
			logger.TagError,
		),
	})
}

func NewResponseLogger(w io.Writer) fiber.Handler {
	return logger.New(logger.Config{
		Next: func(c fiber.Ctx) bool {
			return ctxutil.SkipResponseLog(c)
		},
		Output:     w,
		TimeFormat: timeFormat,
		TimeZone:   timeZone,
		Format: fmt.Sprintf("${%s} [%s] ${%s} ${%s} ${%s} | ${%s} | ${%s%s} | ${%s}\n",
			logger.TagTime,
			"Response Log",
			logger.TagStatus,
			logger.TagLatency,
			logger.TagPath,
			logger.TagResBody,
			// ${locals:requestid}
			logger.TagLocals,
			"${respHeader:X-Request-ID}",
			logger.TagError,
		),
	})
}
