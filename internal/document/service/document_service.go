package service

import (
	"context"
	"io"

	"github.com/suryaadi44/eAD-System/internal/document/dto"
)

type DocumentService interface {
	AddTemplate(ctx context.Context, template dto.TemplateRequest, file io.Reader, fileName string) error
	GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error)
	GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error)

	AddDocument(ctx context.Context, document dto.DocumentRequest, userID string) (string, error)
	GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error)
	GetDocumentStatus(ctx context.Context, documentID string) (*dto.DocumentStatusResponse, error)
	GeneratePDFDocument(ctx context.Context, documentID string) ([]byte, error)
	GetApplicantID(ctx context.Context, documentID string) (*string, error)
	VerifyDocument(ctx context.Context, documentID string, verifierID string) error
	SignDocument(ctx context.Context, documentID string, signerID string) error
}
