package impl

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/template/repository"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"strings"

	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

type TemplateRepositoryImpl struct {
	db *gorm.DB
}

func NewTemplateRepositoryImpl(db *gorm.DB) repository.TemplateRepository {
	return &TemplateRepositoryImpl{
		db: db,
	}
}

func (t *TemplateRepositoryImpl) AddTemplate(ctx context.Context, template *entity.Template) error {
	err := t.db.WithContext(ctx).Create(template).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return utils.ErrDuplicateTemplateName
		}
		return err
	}

	return nil
}

func (t *TemplateRepositoryImpl) GetAllTemplate(ctx context.Context) (*entity.Templates, error) {
	var templates entity.Templates
	err := t.db.WithContext(ctx).
		Preload("Fields").
		Find(&templates).Error
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, utils.ErrTemplateNotFound
	}

	return &templates, nil
}

func (t *TemplateRepositoryImpl) GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error) {
	var template entity.Template
	err := t.db.WithContext(ctx).
		Preload("Fields").
		First(&template, "id = ?", templateId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrTemplateNotFound
		}

		return nil, err
	}

	return &template, nil
}

func (t *TemplateRepositoryImpl) GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error) {
	var templateFields entity.TemplateFields
	err := t.db.WithContext(ctx).Find(&templateFields, "template_id = ?", templateId).Error
	if err != nil {
		return nil, err
	}

	if len(templateFields) == 0 {
		return nil, utils.ErrTemplateFieldNotFound
	}

	return &templateFields, nil
}
