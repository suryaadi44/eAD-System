package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/internal/template/dto"
	"io"
)

type MockTemplateService struct {
	mock.Mock
}

func (m *MockTemplateService) AddTemplate(ctx context.Context, template *dto.TemplateRequest, file io.Reader, fileName string) error {
	args := m.Called(ctx, template, file, fileName)
	return args.Error(0)
}

func (m *MockTemplateService) GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*dto.TemplatesResponse), args.Error(1)
}

func (m *MockTemplateService) GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*dto.TemplateResponse), args.Error(1)
}
