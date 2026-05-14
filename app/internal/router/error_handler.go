package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

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
		c.Logger().Error(err)
	}

	if c.Request().Method == http.MethodHead {
		if err := c.NoContent(code); err != nil {
			c.Logger().Error(err)
		}
		return
	}

	if err := c.JSON(code, echo.Map{"message": message}); err != nil {
		c.Logger().Error(err)
	}
}
