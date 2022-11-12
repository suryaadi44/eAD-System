package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/template/dto"
	"io"
)

type TemplateService interface {
	AddTemplate(ctx context.Context, template *dto.TemplateRequest, file io.Reader, fileName string) error
	GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error)
	GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error)
}
