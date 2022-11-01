package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User, ctx context.Context) error
}
