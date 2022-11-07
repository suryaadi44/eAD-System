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
	GetBriefDocuments(ctx context.Context, limit int, offset int) (*entity.Documents, error)
	GetBriefDocumentsByApplicant(ctx context.Context, applicantID string, limit int, offset int) (*entity.Documents, error)
	GetDocumentStatus(ctx context.Context, documentID string) (*entity.Document, error)
	GetApplicantID(ctx context.Context, documentID string) (*string, error)
	GetDocumentStage(ctx context.Context, documentID string) (*int, error)
	VerifyDocument(ctx context.Context, document *entity.Document) error
	SignDocument(ctx context.Context, document *entity.Document) error
	DeleteDocument(ctx context.Context, documentID string) error
	UpdateDocument(ctx context.Context, document *entity.Document) error
	UpdateDocumentFields(ctx context.Context, documentFields *entity.DocumentFields) error
}
