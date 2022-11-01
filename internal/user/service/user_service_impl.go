package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/internal/user/repository"
)

type (
	PaswordHashFunction interface {
		GenerateFromPassword(password []byte, cost int) ([]byte, error)
		CompareHashAndPassword(hashedPassword, password []byte) error
	}

	UserServiceImpl struct {
		userRepository repository.UserRepository
		passwordHash   PaswordHashFunction
	}
)

func NewUserServiceImpl(userRepository repository.UserRepository, function PaswordHashFunction) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		passwordHash:   function,
	}
}

func (u *UserServiceImpl) SignUpUser(user *dto.UserSignUpRequest, ctx context.Context) error {
	hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	userEntity := user.ToEntity()
	userEntity.ID = uuid.New().String()

	err = u.userRepository.CreateUser(userEntity, ctx)
	if err != nil {
		return err
	}

	return nil
}
