package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"github.com/yamamoto99/go-template/app/internal/api"
	"github.com/yamamoto99/go-template/app/internal/entity"
	"github.com/yamamoto99/go-template/app/internal/handler"
	appmock "github.com/yamamoto99/go-template/app/test/mock/usecase"
)

func TestWelcomeHandler_GetRandomUser(t *testing.T) {
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	user := &entity.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(user, nil)

	h := handler.NewWelcomeHandler(mockUsecase)

	resp, err := h.GetRandomUser(context.Background(), api.GetRandomUserRequestObject{})
	assert.NoError(t, err)

	successResp, ok := resp.(api.GetRandomUser200JSONResponse)
	assert.True(t, ok, "expected GetRandomUser200JSONResponse, got %T", resp)
	assert.Equal(t, user.ID, successResp.Id)
	assert.Equal(t, user.Name, successResp.Name)
	assert.Equal(t, openapi_types.Email(user.Email), successResp.Email)
	mockUsecase.AssertExpectations(t)
}

func TestWelcomeHandler_GetRandomUser_Error(t *testing.T) {
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	expectedErr := errors.New("usecase error")
	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(nil, expectedErr)

	h := handler.NewWelcomeHandler(mockUsecase)

	resp, err := h.GetRandomUser(context.Background(), api.GetRandomUserRequestObject{})
	assert.Nil(t, resp)

	var he *echo.HTTPError
	assert.ErrorAs(t, err, &he)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
	mockUsecase.AssertExpectations(t)
}
