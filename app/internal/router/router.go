package router

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/yamamoto99/go-template/app/internal/api"
	"github.com/yamamoto99/go-template/app/internal/handler"
	"github.com/yamamoto99/go-template/app/internal/logging"
)

func NewRouter(
	wh handler.WelcomeHandler,
) (*echo.Echo, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("load openapi spec: %w", err)
	}
	// Clear servers so the validator skips host matching; otherwise requests
	// from hosts not listed in the spec (e.g. tests, internal callers) get rejected.
	spec.Servers = nil

	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logging.RequestIDMiddleware())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogError:     true,
		LogRequestID: true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("request_id", v.RequestID),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				slog.LogAttrs(c.Request().Context(), slog.LevelError, "request", attrs...)
			} else {
				slog.LogAttrs(c.Request().Context(), slog.LevelInfo, "request", attrs...)
			}
			return nil
		},
	}))
	e.Use(middleware.BodyLimit("1M"))
	e.Use(oapimiddleware.OapiRequestValidator(spec))
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	// }))

	api.RegisterHandlers(e, api.NewStrictHandler(wh, nil))

	return e, nil
}
