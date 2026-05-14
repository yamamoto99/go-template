package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/yamamoto99/go-template/app/internal/logging"
)

func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	logger := logging.FromContext(c.Request().Context())

	code := http.StatusInternalServerError
	message := http.StatusText(http.StatusInternalServerError)

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		if msg, ok := he.Message.(string); ok && msg != "" {
			message = msg
		} else {
			message = http.StatusText(code)
		}
	} else {
		logger.Error("unhandled error", "err", err)
	}

	if c.Request().Method == http.MethodHead {
		if err := c.NoContent(code); err != nil {
			logger.Error("write response", "err", err)
		}
		return
	}

	if err := c.JSON(code, echo.Map{"message": message}); err != nil {
		logger.Error("write response", "err", err)
	}
}
