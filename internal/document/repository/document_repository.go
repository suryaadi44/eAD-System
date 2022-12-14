package repository

import (
	"context"

	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRepository interface {
	AddDocument(ctx context.Context, document *entity.Document) (string, error)
	GetDocument(ctx context.Context, documentID string) (*entity.Document, error)
	GetBriefDocument(ctx context.Context, documentID string) (*entity.Document, error)
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

	AddDocumentRegister(ctx context.Context, register *entity.Register) (uint, error)
}
