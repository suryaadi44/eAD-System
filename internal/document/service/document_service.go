package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
)

type DocumentService interface {
	AddTemplate(ctx context.Context, template dto.TemplateRequest) error
}
