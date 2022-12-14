package impl

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"strings"

	"github.com/suryaadi44/eAD-System/pkg/config"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	db *gorm.DB
}

func NewDocumentRepositoryImpl(db *gorm.DB) repository.DocumentRepository {
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

func (d *DocumentRepositoryImpl) AddDocument(ctx context.Context, document *entity.Document) (string, error) {
	err := d.db.WithContext(ctx).Omit("Register").Create(document).Error
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

func (d *DocumentRepositoryImpl) GetBriefDocument(ctx context.Context, documentID string) (*entity.Document, error) {
	var document entity.Document
	err := d.db.WithContext(ctx).Model(&entity.Document{}).
		Select("id, register_id, description, created_at, applicant_id, template_id, stage_id").
		Preload("Applicant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name")
		}).
		Preload("Template", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Stage").
		Preload("Register").
		Order("created_at desc").
		First(&document, "id = ?", documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrDocumentNotFound
		}

		return nil, err
	}

	return &document, nil
}

func (d *DocumentRepositoryImpl) GetBriefDocuments(ctx context.Context, limit int, offset int) (*entity.Documents, error) {
	var documents entity.Documents
	err := d.db.WithContext(ctx).Model(&entity.Document{}).
		Select("id, register_id, description, created_at, applicant_id, template_id, stage_id").
		Preload("Applicant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name")
		}).
		Preload("Template", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Stage").
		Preload("Register").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&documents).Error
	if err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		return nil, utils.ErrDocumentNotFound
	}

	return &documents, nil
}

func (d *DocumentRepositoryImpl) GetBriefDocumentsByApplicant(ctx context.Context, applicantID string, limit int, offset int) (*entity.Documents, error) {
	var documents entity.Documents
	err := d.db.WithContext(ctx).Model(&entity.Document{}).
		Select("id, register_id, description, created_at, applicant_id, template_id, stage_id").
		Preload("Applicant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name")
		}).
		Preload("Template", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Stage").
		Preload("Register").
		Where("applicant_id = ?", applicantID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&documents).Error
	if err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		return nil, utils.ErrDocumentNotFound
	}

	return &documents, nil
}

func (d *DocumentRepositoryImpl) GetDocumentStatus(ctx context.Context, documentID string) (*entity.Document, error) {
	var document entity.Document
	err := d.db.WithContext(ctx).
		Preload("Stage").
		Preload("Verifier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name, n_ip, position")
		}).
		Preload("Signer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, name, n_ip, position")
		}).
		First(&document, "id = ?", documentID).Error
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
	err := d.db.WithContext(ctx).
		Model(&entity.Document{}).
		Select("applicant_id").
		First(&applicantID, "id = ?", documentID).Error
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
	err := d.db.WithContext(ctx).
		Model(&entity.Document{}).
		Select("stage_id").
		First(&stage, "id = ?", documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrDocumentNotFound
		}

		return nil, err
	}

	return &stage, nil
}

func (d *DocumentRepositoryImpl) VerifyDocument(ctx context.Context, document *entity.Document) error {
	result := d.db.
		WithContext(ctx).
		Model(&entity.Document{}).
		Where("id = ?", document.ID).
		Updates(document)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrDocumentNotFound
	}

	return nil
}

func (d *DocumentRepositoryImpl) SignDocument(ctx context.Context, document *entity.Document) error {
	result := d.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("id = ?", document.ID).
		Select("SignerID", "SignedAt", "StageID").
		Updates(document)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrDocumentNotFound
	}

	return nil
}

func (d *DocumentRepositoryImpl) DeleteDocument(ctx context.Context, documentID string) error {
	result := d.db.WithContext(ctx).
		Select("DocumentField").
		Delete(&entity.Document{}, "id = ?", documentID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrDocumentNotFound
	}

	return nil
}

func (d *DocumentRepositoryImpl) UpdateDocument(ctx context.Context, document *entity.Document) error {
	result := d.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("id = ?", document.ID).
		Updates(document)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrDocumentNotFound
	}

	return nil
}

func (d *DocumentRepositoryImpl) UpdateDocumentFields(ctx context.Context, documentFields *entity.DocumentFields) error {
	for _, documentField := range *documentFields {
		result := d.db.WithContext(ctx).
			Model(&entity.DocumentField{}).
			Where("id = ?", documentField.ID).
			Where("document_id = ?", documentField.DocumentID).
			Updates(documentField)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return utils.ErrFieldNotFound
		}
	}

	return nil
}

func (d *DocumentRepositoryImpl) AddDocumentRegister(ctx context.Context, register *entity.Register) (uint, error) {
	result := d.db.WithContext(ctx).Create(register)
	if result.Error != nil {
		return 0, result.Error
	}

	return register.ID, nil
}
