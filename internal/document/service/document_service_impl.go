package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
	"github.com/suryaadi44/eAD-System/pkg/pdf"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type DocumentServiceImpl struct {
	documentRepository repository.DocumentRepository
	pdfService         pdf.PDFService
}

func NewDocumentServiceImpl(documentRepository repository.DocumentRepository, pdfgService pdf.PDFService) DocumentService {
	return &DocumentServiceImpl{
		documentRepository: documentRepository,
		pdfService:         pdfgService,
	}
}

func (d *DocumentServiceImpl) AddTemplate(ctx context.Context, template dto.TemplateRequest, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	path := filepath.Join("./template", file.Filename)

	// check if file already exist
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file '%s' already exist", file.Filename)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	if err = dst.Close(); err != nil {
		return err
	}

	templateEntity := template.ToEntity()
	templateEntity.Path = path

	err = d.documentRepository.AddTemplate(ctx, templateEntity)
	if err != nil {
		return err
	}

	return nil
}

func (d *DocumentServiceImpl) GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error) {
	templates, err := d.documentRepository.GetAllTemplate(ctx)
	if err != nil {
		return nil, err
	}

	templateResponse := dto.NewTemplatesResponse(templates)

	return templateResponse, nil
}

func (d *DocumentServiceImpl) GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error) {
	template, err := d.documentRepository.GetTemplateDetail(ctx, templateId)
	if err != nil {
		return nil, err
	}

	templateResponse := dto.NewTemplateResponse(template)

	return templateResponse, nil
}

func (d *DocumentServiceImpl) AddDocument(ctx context.Context, document dto.DocumentRequest, userID string) (string, error) {
	keyList, err := d.documentRepository.GetTemplateFields(ctx, document.TemplateID)
	if err != nil {
		return "", err
	}

	// validate document fields with template fields
	for _, key := range *keyList {
		match := false
		for _, field := range document.Fields {
			if key.ID == field.FieldID {
				match = true
				break
			}
		}

		if !match {
			return "", utils.ErrFieldNotMatch
		}
	}

	var documentEntity = document.ToEntity()
	documentEntity.ID = uuid.New().String()
	documentEntity.ApplicantID = userID
	id, err := d.documentRepository.AddDocument(ctx, documentEntity)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (d *DocumentServiceImpl) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	document, err := d.documentRepository.GetDocument(ctx, documentID)
	if err != nil {
		return nil, err
	}

	var documentResponse = dto.NewDocumentResponse(document)

	return documentResponse, nil
}
