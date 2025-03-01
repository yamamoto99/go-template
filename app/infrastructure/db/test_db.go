package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewTestDB はテスト用のデータベース接続を作成します
func NewTestDB() *gorm.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("🔴 Error connecting to test database: %s", err)
	}
	return db
}

// CleanupTestDB はテスト用のデータベースをクリーンアップします
func CleanupTestDB(db *gorm.DB) {
	// テーブルのデータをクリアする
	db.Exec("DELETE FROM users")
}
