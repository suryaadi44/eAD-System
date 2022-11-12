package mock

import "github.com/stretchr/testify/mock"

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
