package repository

import (
	"context"

	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type TemplateRepository interface {
	AddTemplate(ctc context.Context, template *entity.Template) error
	GetAllTemplate(ctx context.Context) (*entity.Templates, error)
	GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error)
	GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error)
}
