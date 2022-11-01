package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRepository interface {
	AddTemplate(ctc context.Context, template *entity.Template) (uint, error)
	AddTemplateFields(ctc context.Context, templateField *entity.TemplateFields) error
}
