package mock

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"html/template"
)

type MockRenderService struct {
	mock.Mock
}

func (m *MockRenderService) GenerateSignature(signer entity.User) (*template.HTML, error) {
	args := m.Called(signer)
	return args.Get(0).(*template.HTML), args.Error(1)
}

func (m *MockRenderService) GenerateFooter(document *entity.Document) (*template.HTML, error) {
	args := m.Called(document)
	return args.Get(0).(*template.HTML), args.Error(1)
}

func (m *MockRenderService) GenerateHTMLDocument(docTemplate *entity.Template, data *map[string]interface{}) (*bytes.Buffer, error) {
	args := m.Called(docTemplate, data)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}
