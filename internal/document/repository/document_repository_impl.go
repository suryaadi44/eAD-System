package repository

import (
	"context"
	"strings"

	"github.com/suryaadi44/eAD-System/pkg/config"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	db *gorm.DB
}

func NewDocumentRepositoryImpl(db *gorm.DB) DocumentRepository {
	documentRepository := &DocumentRepositoryImpl{
		db: db,
	}

	err := documentRepository.InitDefaultStage()
	if err != nil {
		panic(err)
	}

	return documentRepository
}

func (d *DocumentRepositoryImpl) InitDefaultStage() error {
	var count int64
	err := d.db.Model(&entity.Stage{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count != 0 {
		return nil
	}

	for idx, stage := range config.DefaultDocumentStage {
		err := d.db.Create(&entity.Stage{
			ID:     idx + 1,
			Status: stage,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
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

func (d *DocumentRepositoryImpl) GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error) {
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

func (d *DocumentRepositoryImpl) GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error) {
	var templateFields entity.TemplateFields
	err := d.db.WithContext(ctx).Find(&templateFields, "template_id = ?", templateId).Error
	if err != nil {
		return nil, err
	}

	if len(templateFields) == 0 {
		return nil, utils.ErrTemplateFieldNotFound
	}

	return &templateFields, nil
}

func (d *DocumentRepositoryImpl) AddDocument(ctx context.Context, document *entity.Document) (string, error) {
	err := d.db.WithContext(ctx).Create(document).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return "", utils.ErrDuplicateRegister
		}

		return "", err
	}

	return document.ID, nil
}

func (d *DocumentRepositoryImpl) GetDocument(ctx context.Context, documentID string) (*entity.Document, error) {
	var document entity.Document
	err := d.db.WithContext(ctx).
		Preload("Applicant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name")
		}).
		Preload("Verifier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name, n_ip, position")
		}).
		Preload("Signer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name, n_ip, position")
		}).
		Preload("Template").
		Preload("Fields").
		Preload("Stage").
		Preload("Fields.TemplateField").First(&document, "id = ?", documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrDocumentNotFound
		}

		return nil, err
	}

	return &document, nil
}

func (d *DocumentRepositoryImpl) GetApplicantID(ctx context.Context, documentID string) (*string, error) {
	var applicantID string
	err := d.db.WithContext(ctx).Model(&entity.Document{}).Select("applicant_id").First(&applicantID, "id = ?", documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrDocumentNotFound
		}

		return nil, err
	}

	return &applicantID, nil
}

func (d *DocumentRepositoryImpl) GetDocumentStage(ctx context.Context, documentID string) (*int, error) {
	var stage int
	err := d.db.WithContext(ctx).Model(&entity.Document{}).Select("stage_id").First(&stage, "id = ?", documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrDocumentNotFound
		}

		return nil, err
	}

	return &stage, nil
}

func (d *DocumentRepositoryImpl) VerifyDocument(ctx context.Context, document *entity.Document) error {
	result := d.db.WithContext(ctx).Model(&entity.Document{}).Where("id = ?", document.ID).Select("VerifierID", "VerifiedAt", "StageID").Updates(document)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrDocumentNotFound
	}

	return nil
}
