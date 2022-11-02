package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/suryaadi44/eAD-System/pkg/html"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/pdf"
	"github.com/suryaadi44/eAD-System/pkg/utils"
)

type DocumentServiceImpl struct {
	documentRepository repository.DocumentRepository
	pdfService         pdf.PDFService
	renderService      html.RenderService
}

func NewDocumentServiceImpl(documentRepository repository.DocumentRepository, pdfgService pdf.PDFService, renderService html.RenderService) DocumentService {
	return &DocumentServiceImpl{
		documentRepository: documentRepository,
		pdfService:         pdfgService,
		renderService:      renderService,
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
	tmpl, err := d.documentRepository.GetTemplateDetail(ctx, templateId)
	if err != nil {
		return nil, err
	}

	templateResponse := dto.NewTemplateResponse(tmpl)

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

func (d *DocumentServiceImpl) GeneratePDFDocument(ctx context.Context, documentID string) ([]byte, error) {
	document, err := d.documentRepository.GetDocument(ctx, documentID)
	if err != nil {
		return nil, err
	}

	fieldsMap, err := d.fillMapFields(document)
	if err != nil {
		return nil, err
	}

	generatedHTML, err := d.renderService.GenerateHTMLDocument(&document.Template, fieldsMap)
	if err != nil {
		return nil, err
	}

	generatedPDF, err := d.pdfService.GeneratePDF(generatedHTML, document.Template.MarginTop, document.Template.MarginBottom, document.Template.MarginLeft, document.Template.MarginRight)
	if err != nil {
		return nil, err
	}

	return generatedPDF, nil
}

func (d *DocumentServiceImpl) fillMapFields(document *entity.Document) (*map[string]interface{}, error) {
	fieldsMap := dto.NewFieldsMapResponse(&document.Fields)
	fieldsMap["register"] = document.Register

	if document.SignedAt.IsZero() {
		fieldsMap["signedDate"] = ""
		fieldsMap["signature"] = ""
		fieldsMap["footer"] = ""
		return &fieldsMap, nil
	}

	signature, err := d.renderService.GenerateSignature(document.Signer)
	if err != nil {
		return nil, err
	}
	fieldsMap["signature"] = signature

	footer, err := d.renderService.GenerateFooter(document)
	if err != nil {
		return nil, err
	}
	fieldsMap["footer"] = footer

	fieldsMap["signedDate"] = document.SignedAt.Format("02 January 2006")

	return &fieldsMap, nil
}

func (d *DocumentServiceImpl) GetApplicantID(ctx context.Context, documentID string) (*string, error) {
	return d.documentRepository.GetApplicantID(ctx, documentID)
}

func (d *DocumentServiceImpl) VerifyDocument(ctx context.Context, documentID string, verifierID string) error {
	// check stage
	stage, err := d.documentRepository.GetDocumentStage(ctx, documentID)
	if err != nil {
		return err
	}

	if *stage > 1 {
		return utils.ErrAlreadyVerified
	}

	var documentEntity = entity.Document{}
	documentEntity.ID = documentID
	documentEntity.VerifierID = verifierID
	documentEntity.VerifiedAt = time.Now()
	documentEntity.StageID = 2

	return d.documentRepository.VerifyDocument(ctx, &documentEntity)
}
