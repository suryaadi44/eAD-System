package repository

import (
	"context"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	db *gorm.DB
}

func NewDocumentRepositoryImpl(db *gorm.DB) DocumentRepository {
	return &DocumentRepositoryImpl{db}
}

func (d *DocumentRepositoryImpl) AddTemplate(ctx context.Context, template *entity.Template) error {
	result := d.db.WithContext(ctx).Create(template)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *DocumentRepositoryImpl) GetAllTemplate(ctx context.Context) (*entity.Templates, error) {
	var templates entity.Templates
	err := d.db.WithContext(ctx).Preload("Fields").Find(&templates).Error
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, utils.ErrTemplateNotFound
	}

	return &templates, nil
}

func (d *DocumentRepositoryImpl) GetTemplateDetail(ctx context.Context, templateId int64) (*entity.Template, error) {
	var template entity.Template
	err := d.db.WithContext(ctx).Preload("Fields").First(&template, "id = ?", templateId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrTemplateNotFound
		}

		return nil, err
	}

	return &template, nil
}
