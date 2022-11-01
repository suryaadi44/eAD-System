package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	db *gorm.DB
}

func NewDocumentRepositoryImpl(db *gorm.DB) DocumentRepository {
	return &DocumentRepositoryImpl{db}
}

func (d *DocumentRepositoryImpl) AddTemplate(ctx context.Context, template *entity.Template) (uint, error) {
	result := d.db.WithContext(ctx).Create(template)
	if result.Error != nil {
		return 0, result.Error
	}

	return template.ID, nil
}

func (d *DocumentRepositoryImpl) AddTemplateFields(ctx context.Context, templateField *entity.TemplateFields) error {
	err := d.db.WithContext(ctx).Create(templateField).Error
	if err != nil {
		return err
	}

	return nil
}
