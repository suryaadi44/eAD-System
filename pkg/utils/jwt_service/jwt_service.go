package jwt_service

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type JWTService interface {
	GenerateToken(user *entity.User) (string, error)
	GetClaims(c *echo.Context) jwt.MapClaims
}
