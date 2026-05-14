package usecase

import (
	"errors"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/yamamoto99/go-template/app/internal/entity"
	"github.com/yamamoto99/go-template/app/internal/repository"
)

type WelcomeUsecase interface {
	GetRandomUser(c echo.Context) (*entity.User, error)
}

type welcomeUsecase struct {
	wr repository.WelcomeRepository
}

func NewWelcomeUsecase(r repository.WelcomeRepository) WelcomeUsecase {
	return &welcomeUsecase{wr: r}
}

func (u *welcomeUsecase) GetRandomUser(c echo.Context) (*entity.User, error) {
	users, err := u.wr.GetAllUsers(c)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("users not found")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(users))

	return &users[randomIndex], nil
}
