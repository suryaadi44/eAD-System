package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRepository interface {
	AddTemplate(ctc context.Context, template *entity.Template) error
	GetAllTemplate(ctx context.Context) (*entity.Templates, error)
	GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error)
	GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error)

	AddDocument(ctx context.Context, document *entity.Document) (string, error)
	GetDocument(ctx context.Context, documentID string) (*entity.Document, error)
	GetApplicantID(ctx context.Context, documentID string) (*string, error)
}
