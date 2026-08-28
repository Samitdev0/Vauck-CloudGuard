package audit

import (
	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
)

const ContextKey = "audit.context"

type Context struct {
	RequestID string
	Method    string
	Path      string
	IP        string
	UserAgent string

	Identity *auth.Identity
}

func FromEcho(c echo.Context) *Context {
	ctx := &Context{
		RequestID: requestID(c),
		Method:    c.Request().Method,
		Path:      c.Request().URL.Path,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	if identity, ok := c.Get(auth.ContextIdentity).(*auth.Identity); ok {
		ctx.Identity = identity
	}

	return ctx
}

func IsAuthenticated(ctx *Context) bool {
	return ctx != nil && ctx.Identity != nil
}

func requestID(c echo.Context) string {
	if id, ok := c.Get("request_id").(string); ok {
		return id
	}

	return ""
}
