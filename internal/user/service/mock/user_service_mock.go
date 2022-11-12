package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUpUser(ctx context.Context, user *dto.UserSignUpRequest) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserService) LogInUser(ctx context.Context, user *dto.UserLoginRequest) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetBriefUsers(ctx context.Context, page int, limit int) (*dto.BriefUsersResponse, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(*dto.BriefUsersResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID string, request *dto.UserUpdateRequest) error {
	args := m.Called(ctx, userID, request)
	return args.Error(0)
}
