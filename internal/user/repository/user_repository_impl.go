package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			switch {
			case strings.Contains(err.Error(), "username"):
				return utils.ErrUsernameAlreadyExist
			case strings.Contains(err.Error(), "n_ip"):
				return utils.ErrNIPAlreadyExist
			case strings.Contains(err.Error(), "nik"):
				return utils.ErrNIKAlreadyExist
			}
		}

		return err
	}

	return nil
}

func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Select([]string{"id", "username", "password"}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}