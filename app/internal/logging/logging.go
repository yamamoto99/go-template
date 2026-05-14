package logging

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type ctxKey struct{}

var requestIDKey = ctxKey{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func FromContext(ctx context.Context) *slog.Logger {
	if id := RequestIDFromContext(ctx); id != "" {
		return slog.Default().With("request_id", id)
	}
	return slog.Default()
}

// RequestIDMiddleware reads the X-Request-ID header set by echo's RequestID
// middleware and stores the value in the request context so downstream
// loggers can attach it.
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Response().Header().Get(echo.HeaderXRequestID)
			if id != "" {
				req := c.Request()
				c.SetRequest(req.WithContext(WithRequestID(req.Context(), id)))
			}
			return next(c)
		}
	}
}
