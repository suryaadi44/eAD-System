package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"mime/multipart"
)

type DocumentService interface {
	AddTemplate(ctx context.Context, template dto.TemplateRequest, file *multipart.FileHeader) error
	GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error)
	GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error)

	AddDocument(ctx context.Context, document dto.DocumentRequest, userID string) (string, error)
}
