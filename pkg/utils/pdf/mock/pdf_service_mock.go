package mock

import (
	"bytes"
	"github.com/stretchr/testify/mock"
)

type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) GeneratePDF(data *bytes.Buffer, marginTop uint, marginBottom uint, marginLeft uint, marginRight uint) ([]byte, error) {
	args := m.Called(data, marginTop, marginBottom, marginLeft, marginRight)
	return args.Get(0).([]byte), args.Error(1)
}
