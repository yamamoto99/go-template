package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yamamoto99/go-template/app/internal/repository"
	"github.com/yamamoto99/go-template/app/test"
)

func TestWelcomeRepository_GetAllUsers(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	expectedUser := test.SeedTestUser(t, db)
	repo := repository.NewWelcomeRepository(db)

	users, err := repo.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, expectedUser.ID, users[0].ID)
	assert.Equal(t, expectedUser.Name, users[0].Name)
	assert.Equal(t, expectedUser.Email, users[0].Email)
}

func TestWelcomeRepository_GetAllUsers_Empty(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	repo := repository.NewWelcomeRepository(db)

	users, err := repo.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
