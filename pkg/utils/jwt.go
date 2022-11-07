package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"time"
)

type JWTService struct {
	secretKey string
	exp       time.Duration
}

func NewJWTService(secretKey string, exp time.Duration) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		exp:       exp,
	}
}

func (j *JWTService) GenerateToken(user *entity.User) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(j.exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
