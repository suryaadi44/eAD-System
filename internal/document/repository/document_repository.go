package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRepository interface {
	AddTemplate(ctc context.Context, template *entity.Template) error
	GetAllTemplate(ctx context.Context) (*entity.Templates, error)
}
