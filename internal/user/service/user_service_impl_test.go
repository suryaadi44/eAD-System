package service

import (
	"context"
	"errors"
	"testing"

	"github.com/suryaadi44/eAD-System/pkg/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetBriefUsers(ctx context.Context, limit int, offset int) (*entity.Users, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).(*entity.Users), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockPasswordHashFunction struct {
	mock.Mock
}

func (m *MockPasswordHashFunction) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPasswordHashFunction) CompareHashAndPassword(hashedPassword, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

type TestSuiteUserService struct {
	suite.Suite
	mockUserRepository *MockUserRepository
	mockPasswordHash   *MockPasswordHashFunction
	mockJWTService     *MockJWTService
	userService        UserService
}

func (t *TestSuiteUserService) SetupTest() {
	t.mockUserRepository = new(MockUserRepository)
	t.mockPasswordHash = new(MockPasswordHashFunction)
	t.mockJWTService = new(MockJWTService)
	t.userService = NewUserServiceImpl(t.mockUserRepository, t.mockPasswordHash, t.mockJWTService)
}

func (t *TestSuiteUserService) TearDownTest() {
	t.mockUserRepository = nil
	t.mockPasswordHash = nil
	t.mockJWTService = nil
	t.userService = nil
}

func (t *TestSuiteUserService) TestSignUpUser_Success() {
	t.mockPasswordHash.On("GenerateFromPassword", []byte("password"), 10).Return([]byte("hashedPassword"), nil)
	t.mockUserRepository.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	err := t.userService.SignUpUser(context.Background(), &dto.UserSignUpRequest{
		Username: "username",
		Password: "password",
	})

	t.NoError(err)
}

func (t *TestSuiteUserService) TestSignUpUser_FailedHashing() {
	t.mockPasswordHash.On("GenerateFromPassword", []byte("password"), 10).Return([]byte(""), errors.New("error"))

	err := t.userService.SignUpUser(context.Background(), &dto.UserSignUpRequest{
		Username: "username",
		Password: "password",
	})

	t.Error(errors.New("error"), err)
}

func (t *TestSuiteUserService) TestSignUpUser_FailedCreateUser() {
	t.mockPasswordHash.On("GenerateFromPassword", []byte("password"), 10).Return([]byte("hashedPassword"), nil)
	t.mockUserRepository.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := t.userService.SignUpUser(context.Background(), &dto.UserSignUpRequest{
		Username: "username",
		Password: "password",
	})

	t.Error(errors.New("error"), err)
}

func (t *TestSuiteUserService) TestLoginUser_Success() {
	t.mockUserRepository.On("FindByUsername", mock.Anything, "username").Return(&entity.User{
		Username: "username",
		Password: "hashedPassword",
	}, nil)
	t.mockPasswordHash.On("CompareHashAndPassword", []byte("hashedPassword"), []byte("password")).Return(nil)
	t.mockJWTService.On("GenerateToken", mock.Anything).Return("token", nil)

	resp, err := t.userService.LogInUser(context.Background(), &dto.UserLoginRequest{
		Username: "username",
		Password: "password",
	})

	t.NoError(err)
	t.Equal("token", resp)
}

func (t *TestSuiteUserService) TestLoginUser_FailedFindByUsername() {
	t.mockUserRepository.On("FindByUsername", mock.Anything, "username").Return(&entity.User{}, errors.New("error"))

	_, err := t.userService.LogInUser(context.Background(), &dto.UserLoginRequest{
		Username: "username",
		Password: "password",
	})

	t.Error(err)
}

func (t *TestSuiteUserService) TestLoginUser_FailedUserNotFound() {
	t.mockUserRepository.On("FindByUsername", mock.Anything, "username").Return(&entity.User{}, utils.ErrUserNotFound)

	_, err := t.userService.LogInUser(context.Background(), &dto.UserLoginRequest{
		Username: "username",
		Password: "password",
	})

	t.Equal(utils.ErrInvalidCredentials, err)
}

func (t *TestSuiteUserService) TestLoginUser_FailedCompareHashAndPassword() {
	t.mockUserRepository.On("FindByUsername", mock.Anything, "username").Return(&entity.User{
		Username: "username",
		Password: "hashedPassword",
	}, nil)
	t.mockPasswordHash.On("CompareHashAndPassword", []byte("hashedPassword"), []byte("password")).Return(errors.New("error"))

	_, err := t.userService.LogInUser(context.Background(), &dto.UserLoginRequest{
		Username: "username",
		Password: "password",
	})

	t.Equal(utils.ErrInvalidCredentials, err)
}

func (t *TestSuiteUserService) TestLoginUser_FailedGenerateToken() {
	t.mockUserRepository.On("FindByUsername", mock.Anything, "username").Return(&entity.User{
		Username: "username",
		Password: "hashedPassword",
	}, nil)
	t.mockPasswordHash.On("CompareHashAndPassword", []byte("hashedPassword"), []byte("password")).Return(nil)
	t.mockJWTService.On("GenerateToken", mock.Anything).Return("", errors.New("error"))

	_, err := t.userService.LogInUser(context.Background(), &dto.UserLoginRequest{
		Username: "username",
		Password: "password",
	})

	t.Error(err)
}

func (t *TestSuiteUserService) TestGetBriefUsers_Success() {
	t.mockUserRepository.On("GetBriefUsers", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Users{
		{
			Username: "username1",
		},
	}, nil)

	resp, err := t.userService.GetBriefUsers(context.Background(), 1, 1)

	t.NoError(err)
	t.Equal(&dto.BriefUsersResponse{
		{
			Username: "username1",
		},
	}, resp)
}

func (t *TestSuiteUserService) TestGetBriefUsers_RepoError() {
	t.mockUserRepository.On("GetBriefUsers", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Users{}, errors.New("error"))

	_, err := t.userService.GetBriefUsers(context.Background(), 1, 1)

	t.Error(err)
}

func (t *TestSuiteUserService) TestUpdateUser_Success() {
	t.mockUserRepository.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	err := t.userService.UpdateUser(context.Background(), "userid", &dto.UserUpdateRequest{
		Username: "username",
	})

	t.NoError(err)
}

func (t *TestSuiteUserService) TestUpdateUser_SuccessWithPassword() {
	t.mockUserRepository.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
	t.mockPasswordHash.On("GenerateFromPassword", mock.Anything, 10).Return([]byte("hashedPassword"), nil)

	err := t.userService.UpdateUser(context.Background(), "userid", &dto.UserUpdateRequest{
		Username: "username",
		Password: "password",
	})

	t.NoError(err)
}

func (t *TestSuiteUserService) TestUpdateUser_PasswordHashError() {
	t.mockUserRepository.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
	t.mockPasswordHash.On("GenerateFromPassword", mock.Anything, 10).Return(([]byte)(nil), errors.New("error"))

	err := t.userService.UpdateUser(context.Background(), "userid", &dto.UserUpdateRequest{
		Username: "username",
		Password: "password",
	})

	t.Error(err)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestSuiteUserService))
}
