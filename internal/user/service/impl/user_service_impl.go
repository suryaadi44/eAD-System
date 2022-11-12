package impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/internal/user/repository"
	"github.com/suryaadi44/eAD-System/internal/user/service"
	error2 "github.com/suryaadi44/eAD-System/pkg/utils"
	"github.com/suryaadi44/eAD-System/pkg/utils/jwt_service"
	"github.com/suryaadi44/eAD-System/pkg/utils/password"
)

type (
	UserServiceImpl struct {
		userRepository repository.UserRepository
		passwordHash   password.PasswordFunc
		jwtService     jwt_service.JWTService
	}
)

func NewUserServiceImpl(userRepository repository.UserRepository, function password.PasswordFunc, jwt jwt_service.JWTService) service.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		passwordHash:   function,
		jwtService:     jwt,
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
	userEntity.Role = 1

	err = u.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) LogInUser(ctx context.Context, user *dto.UserLoginRequest) (string, error) {
	userEntity, err := u.userRepository.FindByUsername(ctx, user.Username)
	if err != nil {
		if err == error2.ErrUserNotFound {
			return "", error2.ErrInvalidCredentials
		}

		return "", err
	}

	err = u.passwordHash.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil {
		return "", error2.ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(userEntity)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserServiceImpl) GetBriefUsers(ctx context.Context, page int, limit int) (*dto.BriefUsersResponse, error) {
	offset := (page - 1) * limit

	users, err := u.userRepository.GetBriefUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.NewBriefUsersResponse(users), nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, userID string, request *dto.UserUpdateRequest) error {
	user := request.ToEntity()
	user.ID = userID

	if user.Password != "" {
		hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return u.userRepository.UpdateUser(ctx, user)
}
