package utils

import "golang.org/x/crypto/bcrypt"

type PasswordFunc struct {
}

func (p PasswordFunc) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (p PasswordFunc) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
