package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	GetAllUser(ctx context.Context) (*entity.Users, error)
}
