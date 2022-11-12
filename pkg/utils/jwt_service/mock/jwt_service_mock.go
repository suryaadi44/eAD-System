package mock

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GetClaims(c *echo.Context) jwt.MapClaims {
	args := m.Called(c)
	return args.Get(0).(jwt.MapClaims)
}
