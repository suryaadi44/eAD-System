package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/config"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	userRepository := &UserRepositoryImpl{
		db: db,
	}

	err := userRepository.InitDefaultUser()
	if err != nil {
		panic(err)
	}

	return userRepository
}

func (u *UserRepositoryImpl) InitDefaultUser() error {
	var count int64
	err := u.db.Model(&entity.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count != 0 {
		return nil
	}

	err = u.db.Create(config.DefaultUser).Error
	if err != nil {
		return err
	}

	return nil
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
	err := u.db.WithContext(ctx).Select([]string{"id", "username", "password", "role"}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *UserRepositoryImpl) GetAllUser(ctx context.Context) (*entity.Users, error) {
	var users entity.Users
	err := u.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, utils.ErrUserNotFound
	}

	return &users, nil
}
