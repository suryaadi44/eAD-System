package utils

import "golang.org/x/crypto/bcrypt"

type PasswordFunc struct {
}

func (PasswordFunc) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (PasswordFunc) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
