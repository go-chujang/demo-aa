package ctxutil

import "github.com/gofiber/fiber/v3"

const (
	onlyOwner       = "_routechart_only_owner"
	skipAuth        = "_routechart_skip_auth"
	skipRequestLog  = "_routechart_skip_request_log"
	skipResponseLog = "_routechart_skip_response_log"
)

// getter
func OnlyOwner(c fiber.Ctx) bool       { return MustLocals[string, bool](c, onlyOwner) }
func SkipAuth(c fiber.Ctx) bool        { return MustLocals[string, bool](c, skipAuth) }
func SkipRequestLog(c fiber.Ctx) bool  { return MustLocals[string, bool](c, skipRequestLog) }
func SkipResponseLog(c fiber.Ctx) bool { return MustLocals[string, bool](c, skipResponseLog) }

type Policy struct {
	onlyOwner       bool
	skipAuth        bool
	skipRequestLog  bool
	skipResponseLog bool
}

func (p Policy) Stores(c fiber.Ctx) {
	if p.onlyOwner {
		c.Locals(onlyOwner, true)
	}
	if p.skipAuth {
		c.Locals(skipAuth, true)
	}
	if p.skipRequestLog {
		c.Locals(skipRequestLog, true)
	}
	if p.skipResponseLog {
		c.Locals(skipResponseLog, true)
	}
}

func (p Policy) IsOnlyOwner() bool       { return p.onlyOwner }
func (p Policy) IsSkipAuth() bool        { return p.skipAuth }
func (p Policy) IsSkipRequestLog() bool  { return p.skipRequestLog }
func (p Policy) IsSkipResponseLog() bool { return p.skipResponseLog }

func NewPolicy() *Policy {
	return &Policy{
		onlyOwner:       false,
		skipAuth:        false,
		skipRequestLog:  false,
		skipResponseLog: false,
	}
}

func (p *Policy) OnlyOwner() *Policy {
	p.onlyOwner = true
	return p
}

func (p *Policy) SkipAll() *Policy {
	p.skipAuth = true
	p.skipRequestLog = true
	p.skipResponseLog = true
	return p
}

func (p *Policy) SkipAuth() *Policy {
	p.skipAuth = true
	return p
}

func (p *Policy) SkipLog() *Policy {
	p.skipRequestLog = true
	p.skipResponseLog = true
	return p
}

func (p *Policy) SkipRequestLog() *Policy {
	p.skipRequestLog = true
	return p
}

func (p *Policy) SkipResponseLog() *Policy {
	p.skipResponseLog = true
	return p
}
