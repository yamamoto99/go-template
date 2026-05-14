package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/yamamoto99/go-template/app/internal/usecase"
)

type WelcomeHandler interface {
	GetRandomUser(c echo.Context) error
}

type welcomeHandler struct {
	wu usecase.WelcomeUsecase
}

func NewWelcomeHandler(u usecase.WelcomeUsecase) WelcomeHandler {
	return &welcomeHandler{wu: u}
}

func (h *welcomeHandler) GetRandomUser(c echo.Context) error {
	user, err := h.wu.GetRandomUser(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, user)
}
