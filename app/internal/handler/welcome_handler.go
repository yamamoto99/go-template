package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/yamamoto99/go-template/app/internal/api"
	"github.com/yamamoto99/go-template/app/internal/logging"
	"github.com/yamamoto99/go-template/app/internal/usecase"
)

type WelcomeHandler interface {
	GetRandomUser(ctx context.Context, request api.GetRandomUserRequestObject) (api.GetRandomUserResponseObject, error)
}

type welcomeHandler struct {
	wu usecase.WelcomeUsecase
}

func NewWelcomeHandler(u usecase.WelcomeUsecase) WelcomeHandler {
	return &welcomeHandler{wu: u}
}

func (h *welcomeHandler) GetRandomUser(ctx context.Context, _ api.GetRandomUserRequestObject) (api.GetRandomUserResponseObject, error) {
	user, err := h.wu.GetRandomUser(ctx)
	if err != nil {
		logging.FromContext(ctx).Error("get random user", "err", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return api.GetRandomUser200JSONResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
