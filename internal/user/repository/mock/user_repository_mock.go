package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetBriefUsers(ctx context.Context, limit int, offset int) (*entity.Users, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).(*entity.Users), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
