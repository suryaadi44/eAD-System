package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/internal/user/repository"
)

type (
	PasswordHashFunction interface {
		GenerateFromPassword(password []byte, cost int) ([]byte, error)
		CompareHashAndPassword(hashedPassword, password []byte) error
	}

	UserServiceImpl struct {
		userRepository repository.UserRepository
		passwordHash   PasswordHashFunction
	}
)

func NewUserServiceImpl(userRepository repository.UserRepository, function PasswordHashFunction) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		passwordHash:   function,
	}
}

func (u *UserServiceImpl) SignUpUser(ctx context.Context, user *dto.UserSignUpRequest) error {
	hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	userEntity := user.ToEntity()
	userEntity.ID = uuid.New().String()

	err = u.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		return err
	}

	return nil
}
