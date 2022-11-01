package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
}
