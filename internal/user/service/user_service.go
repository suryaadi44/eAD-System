package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
)

type UserService interface {
	SignUpUser(ctx context.Context, user *dto.UserSignUpRequest) error
	LogInUser(ctx context.Context, user *dto.UserLoginRequest) (string, error)
}
