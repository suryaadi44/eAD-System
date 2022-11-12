package password

import "golang.org/x/crypto/bcrypt"

type PasswordFuncImpl struct {
}

func NewPasswordFuncImpl() PasswordFunc {
	return &PasswordFuncImpl{}
}

func (PasswordFuncImpl) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (PasswordFuncImpl) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
