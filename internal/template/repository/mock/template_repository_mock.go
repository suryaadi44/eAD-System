package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type MockTemplateRepository struct {
	mock.Mock
}

func (m *MockTemplateRepository) AddTemplate(ctc context.Context, template *entity.Template) error {
	args := m.Called(ctc, template)
	return args.Error(0)
}

func (m *MockTemplateRepository) GetAllTemplate(ctx context.Context) (*entity.Templates, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.Templates), args.Error(1)
}

func (m *MockTemplateRepository) GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.Template), args.Error(1)
}

func (m *MockTemplateRepository) GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.TemplateFields), args.Error(1)
}
