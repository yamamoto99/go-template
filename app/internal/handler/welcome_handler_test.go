package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"tmp/app/internal/entity"
	"tmp/app/internal/handler"
	appmock "tmp/app/test/mock"
)

func TestWelcomeHandler_GetRandomUser(t *testing.T) {
	// モックの準備
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	// テストデータ
	user := &entity.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// モックの振る舞いを設定
	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(user, nil)

	// テスト対象のハンドラーを初期化
	h := handler.NewWelcomeHandler(mockUsecase)

	// HTTPリクエストのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/random-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	err := h.GetRandomUser(c)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// レスポンスボディの検証
	var responseUser entity.User
	err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, user.Name, responseUser.Name)
	assert.Equal(t, user.Email, responseUser.Email)

	// モックが期待通り呼ばれたことを確認
	mockUsecase.AssertExpectations(t)
}

func TestWelcomeHandler_GetRandomUser_Error(t *testing.T) {
	// モックの準備
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	// モックの振る舞いを設定
	expectedErr := errors.New("usecase error")
	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(nil, expectedErr)

	// テスト対象のハンドラーを初期化
	h := handler.NewWelcomeHandler(mockUsecase)

	// HTTPリクエストのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/random-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	err := h.GetRandomUser(c)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// レスポンスボディの検証
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedErr.Error(), response["error"])

	// モックが期待通り呼ばれたことを確認
	mockUsecase.AssertExpectations(t)
}
