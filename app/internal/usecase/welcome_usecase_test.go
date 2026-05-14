package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"github.com/yamamoto99/go-template/app/internal/entity"
	"github.com/yamamoto99/go-template/app/internal/usecase"
	appmock "github.com/yamamoto99/go-template/app/test/mock/repository"
)

func TestWelcomeUsecase_GetRandomUser(t *testing.T) {
	mockRepo := new(appmock.WelcomeRepositoryMock)

	users := []entity.User{
		{
			ID:        1,
			Name:      "Test User 1",
			Email:     "test1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Test User 2",
			Email:     "test2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := context.Background()
	mockRepo.On("GetAllUsers", testifymock.Anything).Return(users, nil)

	uc := usecase.NewWelcomeUsecase(mockRepo)

	user, err := uc.GetRandomUser(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Contains(t, []int{1, 2}, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestWelcomeUsecase_GetRandomUser_EmptyUsers(t *testing.T) {
	mockRepo := new(appmock.WelcomeRepositoryMock)

	ctx := context.Background()
	mockRepo.On("GetAllUsers", testifymock.Anything).Return([]entity.User{}, nil)

	uc := usecase.NewWelcomeUsecase(mockRepo)

	user, err := uc.GetRandomUser(ctx)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "users not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestWelcomeUsecase_GetRandomUser_Error(t *testing.T) {
	mockRepo := new(appmock.WelcomeRepositoryMock)

	ctx := context.Background()
	expectedErr := errors.New("database error")
	mockRepo.On("GetAllUsers", testifymock.Anything).Return([]entity.User{}, expectedErr)

	uc := usecase.NewWelcomeUsecase(mockRepo)

	user, err := uc.GetRandomUser(ctx)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
