package test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/yamamoto99/go-template/app/infrastructure/config"
	"github.com/yamamoto99/go-template/app/infrastructure/db"
	"github.com/yamamoto99/go-template/app/internal/entity"
)

func loadEnvFile(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Error loading .env file: %v", err)
	}
}

func SetupTestDB(t *testing.T) *gorm.DB {
	loadEnvFile(t)

	cfg, err := config.LoadTest()
	if err != nil {
		t.Fatalf("Error loading test config: %v", err)
	}

	dbConn, err := db.NewTestDB(cfg.TestDB)
	if err != nil {
		t.Fatalf("Error connecting to test database: %v", err)
	}
	tx := dbConn.Begin()
	if tx.Error != nil {
		t.Fatalf("Error beginning test transaction: %v", tx.Error)
	}
	return tx
}

func CleanupDB(t *testing.T, dbConn *gorm.DB) {
	if err := dbConn.Rollback().Error; err != nil {
		t.Fatalf("Error rolling back test transaction: %v", err)
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		t.Fatalf("Error getting DB instance: %v", err)
	}
	sqlDB.Close()
}

func SeedTestUser(t *testing.T, dbConn *gorm.DB) entity.User {
	user := entity.User{
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
