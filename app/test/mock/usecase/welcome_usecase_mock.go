package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/yamamoto99/go-template/app/internal/entity"
)

type WelcomeUsecaseMock struct {
	mock.Mock
}

func (m *WelcomeUsecaseMock) GetRandomUser(ctx context.Context) (*entity.User, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}
