package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
)

type UserService interface {
	SignUpUser(user *dto.UserSignUpRequest, ctx context.Context) error
}
