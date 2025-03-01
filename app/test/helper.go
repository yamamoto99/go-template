package test

import (
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"

	"tmp/app/infrastructure/db"
	"tmp/app/infrastructure/migrate"
	"tmp/app/internal/entity"

	"github.com/joho/godotenv"
)

// loadEnvFile は環境変数ファイルを読み込みます
func loadEnvFile(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

// SetupTestDB はテスト用のDBをセットアップします
func SetupTestDB(t *testing.T) *gorm.DB {
	loadEnvFile(t)

	dbConn := db.NewTestDB()
	migrate.RunMigrations(dbConn)
	db.CleanupTestDB(dbConn)
	return dbConn
}

// CleanupDB はテスト後にDBをクリーンアップします
func CleanupDB(t *testing.T, dbConn *gorm.DB) {
	db.CleanupTestDB(dbConn)
	sqlDB, err := dbConn.DB()
	if err != nil {
		t.Fatalf("Error getting DB instance: %v", err)
	}
	sqlDB.Close()
}

// SeedTestUser はテスト用のユーザーデータを作成します
func SeedTestUser(t *testing.T, dbConn *gorm.DB) entity.User {
	user := entity.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := dbConn.Create(&user).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	return user
}

// SetupEnv はテスト環境変数をセットアップします
func SetupEnv(t *testing.T) {
	loadEnvFile(t)

	// 環境変数が正しく設定されているか確認
	requiredEnvVars := []string{
		"TEST_DB_USER", "TEST_DB_PASSWORD", "TEST_DB_HOST",
		"TEST_DB_PORT", "TEST_DB_NAME",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Required environment variable %s is not set", envVar)
		}
	}
}
