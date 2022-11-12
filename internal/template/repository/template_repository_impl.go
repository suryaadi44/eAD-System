package repository

import (
	"context"
	error2 "github.com/suryaadi44/eAD-System/pkg/utils"
	"strings"

	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

type TemplateRepositoryImpl struct {
	db *gorm.DB
}

func NewTemplateRepositoryImpl(db *gorm.DB) TemplateRepository {
	return &TemplateRepositoryImpl{
		db: db,
	}
}

func (t *TemplateRepositoryImpl) AddTemplate(ctx context.Context, template *entity.Template) error {
	err := t.db.WithContext(ctx).Create(template).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return error2.ErrDuplicateTemplateName
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
		return nil, error2.ErrTemplateNotFound
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
			return nil, error2.ErrTemplateNotFound
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
		return nil, error2.ErrTemplateFieldNotFound
	}

	return &templateFields, nil
}
