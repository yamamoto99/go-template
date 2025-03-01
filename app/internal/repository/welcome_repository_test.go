package repository_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"tmp/app/internal/repository"
	"tmp/app/test"
)

func TestWelcomeRepository_GetAllUsers(t *testing.T) {
	// テスト用DBのセットアップ
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	// テストデータの作成
	expectedUser := test.SeedTestUser(t, db)

	// リポジトリの初期化
	repo := repository.NewWelcomeRepository(db)

	// テスト実行
	e := echo.New()
	ctx := e.NewContext(nil, nil)
	users, err := repo.GetAllUsers(ctx)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, expectedUser.ID, users[0].ID)
	assert.Equal(t, expectedUser.Name, users[0].Name)
	assert.Equal(t, expectedUser.Email, users[0].Email)
}

func TestWelcomeRepository_GetAllUsers_Empty(t *testing.T) {
	// テスト用DBのセットアップ
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	// リポジトリの初期化
	repo := repository.NewWelcomeRepository(db)

	// テスト実行
	e := echo.New()
	ctx := e.NewContext(nil, nil)
	users, err := repo.GetAllUsers(ctx)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
