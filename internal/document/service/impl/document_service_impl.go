package impl

import (
	"context"
	"fmt"
	"github.com/suryaadi44/eAD-System/internal/document/service"
	error2 "github.com/suryaadi44/eAD-System/pkg/utils"
	"time"

	"github.com/suryaadi44/eAD-System/pkg/utils/html"
	"github.com/suryaadi44/eAD-System/pkg/utils/pdf"

	"github.com/google/uuid"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
	repository2 "github.com/suryaadi44/eAD-System/internal/template/repository"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentServiceImpl struct {
	documentRepository repository.DocumentRepository
	templateRepository repository2.TemplateRepository
	pdfService         pdf.PDFService
	renderService      html.RenderService
}

func NewDocumentServiceImpl(documentRepository repository.DocumentRepository, templateRepository repository2.TemplateRepository, pdfgService pdf.PDFService, renderService html.RenderService) service.DocumentService {
	return &DocumentServiceImpl{
		documentRepository: documentRepository,
		templateRepository: templateRepository,
		pdfService:         pdfgService,
		renderService:      renderService,
	}
}

func (d *DocumentServiceImpl) AddDocument(ctx context.Context, document *dto.DocumentRequest, userID string) (string, error) {
	keyList, err := d.templateRepository.GetTemplateFields(ctx, document.TemplateID)
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
			return "", error2.ErrFieldNotMatch
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

func (d *DocumentServiceImpl) GetBriefDocuments(ctx context.Context, applicantID string, role int, page int, limit int) (*dto.BriefDocumentsResponse, error) {
	offset := (page - 1) * limit

	var documents *entity.Documents
	var err error

	if role == 1 {
		documents, err = d.documentRepository.GetBriefDocumentsByApplicant(ctx, applicantID, limit, offset)
	} else {
		documents, err = d.documentRepository.GetBriefDocuments(ctx, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	var response = dto.NewBriefDocumentsResponse(documents)

	return response, nil
}

func (d *DocumentServiceImpl) GetDocumentStatus(ctx context.Context, documentID string) (*dto.DocumentStatusResponse, error) {
	document, err := d.documentRepository.GetDocumentStatus(ctx, documentID)
	if err != nil {
		return nil, err
	}

	var documentStatusResponse = dto.NewDocumentStatusResponse(document)

	return documentStatusResponse, nil
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
	fieldsMap["register"] = document.RegisterID

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

func (d *DocumentServiceImpl) VerifyDocument(ctx context.Context, documentID string, verifierID string, verifyRequest *dto.VerifyDocumentRequest) error {
	briefDocument, err := d.documentRepository.GetBriefDocument(ctx, documentID)
	if err != nil {
		return err
	}

	if briefDocument.StageID > 1 {
		return error2.ErrAlreadyVerified
	}

	var documentEntity = entity.Document{}
	var description string

	if briefDocument.Description == "" {
		if verifyRequest.Description == "" {
			documentEntity.Description = fmt.Sprintf("%s a.n %s", briefDocument.Template.Name, briefDocument.Applicant.Name)
		} else {
			documentEntity.Description = verifyRequest.Description
		}
		description = documentEntity.Description
	} else {
		description = briefDocument.Description
	}

	if briefDocument.RegisterID == 0 {
		if verifyRequest.RegisterID == 0 {
			registerID, err := d.documentRepository.AddDocumentRegister(ctx, &entity.Register{
				Description: description,
			})

			if err != nil {
				return err
			}
			documentEntity.RegisterID = registerID
		} else {
			documentEntity.RegisterID = verifyRequest.RegisterID
		}
	}

	documentEntity.ID = documentID
	documentEntity.VerifierID = verifierID
	documentEntity.VerifiedAt = time.Now()
	documentEntity.StageID = 2

	return d.documentRepository.VerifyDocument(ctx, &documentEntity)
}

func (d *DocumentServiceImpl) SignDocument(ctx context.Context, documentID string, signerID string) error {
	// check stage
	stage, err := d.documentRepository.GetDocumentStage(ctx, documentID)
	if err != nil {
		return err
	}

	if *stage > 2 {
		return error2.ErrAlreadySigned
	} else if *stage < 2 {
		return error2.ErrNotVerifiedYet
	}

	var documentEntity = entity.Document{}
	documentEntity.ID = documentID
	documentEntity.SignerID = signerID
	documentEntity.SignedAt = time.Now()
	documentEntity.StageID = 3

	return d.documentRepository.SignDocument(ctx, &documentEntity)
}

func (d *DocumentServiceImpl) DeleteDocument(ctx context.Context, userID string, role int, documentID string) error {
	if role == 1 {
		applicantID, err := d.documentRepository.GetApplicantID(ctx, documentID)
		if err != nil {
			return err
		}

		if *applicantID != userID {
			return error2.ErrDidntHavePermission
		}
	}

	stage, err := d.documentRepository.GetDocumentStage(ctx, documentID)
	if err != nil {
		return err
	}

	if *stage == 3 {
		return error2.ErrAlreadySigned
	}

	return d.documentRepository.DeleteDocument(ctx, documentID)
}

func (d *DocumentServiceImpl) UpdateDocument(ctx context.Context, document *dto.DocumentUpdateRequest, documentID string) error {
	stage, err := d.documentRepository.GetDocumentStage(ctx, documentID)
	if err != nil {
		return err
	}

	if *stage == 2 {
		return error2.ErrAlreadyVerified
	}

	if *stage == 3 {
		return error2.ErrAlreadySigned
	}

	documentEntity := document.ToEntity()
	documentEntity.ID = documentID

	return d.documentRepository.UpdateDocument(ctx, documentEntity)
}

func (d *DocumentServiceImpl) UpdateDocumentFields(ctx context.Context, userID string, role int, documentID string, fields *dto.FieldsUpdateRequest) error {
	if role == 1 {
		applicantID, err := d.documentRepository.GetApplicantID(ctx, documentID)
		if err != nil {
			return err
		}

		if *applicantID != userID {
			return error2.ErrDidntHavePermission
		}
	}

	stage, err := d.documentRepository.GetDocumentStage(ctx, documentID)
	if err != nil {
		return err
	}

	if *stage == 2 {
		return error2.ErrAlreadyVerified
	}

	if *stage == 3 {
		return error2.ErrAlreadySigned
	}

	fieldsEntity := fields.ToEntity(documentID)
	return d.documentRepository.UpdateDocumentFields(ctx, fieldsEntity)
}
