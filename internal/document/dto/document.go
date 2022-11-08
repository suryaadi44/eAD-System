package dto

import (
	"gorm.io/gorm"
	"time"

	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRequest struct {
	TemplateID uint          `json:"template_id" validate:"required"`
	Fields     FieldsRequest `json:"fields" validate:"required"`
}

type FieldRequest struct {
	FieldID uint   `json:"field_id" validate:"required"`
	Value   string `json:"value" validate:"required"`
}

type FieldsRequest []FieldRequest

func (d *DocumentRequest) ToEntity() *entity.Document {
	var fields entity.DocumentFields
	for _, field := range d.Fields {
		fields = append(fields, entity.DocumentField{
			TemplateFieldID: field.FieldID,
			Value:           field.Value,
		})
	}

	return &entity.Document{
		TemplateID: d.TemplateID,
		Fields:     fields,
	}
}

type DocumentResponse struct {
	ID          string                `json:"id"`
	RegisterID  uint                  `json:"register"`
	Description string                `json:"description"`
	Applicant   dto.ApplicantResponse `json:"applicant"`
	Template    TemplateResponse      `json:"template"`
	Fields      FieldsResponse        `json:"fields"`
	Stage       string                `json:"stage"`
	Verifier    dto.EmployeeResponse  `json:"verifier"`
	VerifiedAt  time.Time             `json:"verified_at"`
	Signer      dto.EmployeeResponse  `json:"signer"`
	SignedAt    time.Time             `json:"signed_at"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

func NewDocumentResponse(document *entity.Document) *DocumentResponse {
	return &DocumentResponse{
		ID:          document.ID,
		RegisterID:  document.RegisterID,
		Description: document.Description,
		Applicant:   *dto.NewApplicantResponse(&document.Applicant),
		Template:    *NewTemplateResponse(&document.Template),
		Fields:      *NewFieldsResponse(&document.Fields),
		Stage:       document.Stage.Status,
		Verifier:    *dto.NewEmployeeResponse(&document.Verifier),
		VerifiedAt:  document.VerifiedAt,
		Signer:      *dto.NewEmployeeResponse(&document.Signer),
		SignedAt:    document.SignedAt,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}
}

type FieldResponse struct {
	ID    uint   `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewFieldResponse(fields *entity.DocumentField) *FieldResponse {
	return &FieldResponse{
		ID:    fields.ID,
		Key:   fields.TemplateField.Key,
		Value: fields.Value,
	}
}

type FieldsResponse []FieldResponse

func NewFieldsResponse(fields *entity.DocumentFields) *FieldsResponse {
	var fieldsResponse FieldsResponse
	for _, field := range *fields {
		fieldsResponse = append(fieldsResponse, *NewFieldResponse(&field))
	}

	return &fieldsResponse
}

func NewFieldsMapResponse(fields *entity.DocumentFields) map[string]interface{} {
	var fieldsResponse = make(map[string]interface{})
	for _, field := range *fields {
		fieldsResponse[field.TemplateField.Key] = field.Value
	}

	return fieldsResponse
}

type DocumentStatusResponse struct {
	ID          string               `json:"id"`
	Description string               `json:"description"`
	RegisterID  uint                 `json:"register"`
	Stage       string               `json:"stage"`
	Verifier    dto.EmployeeResponse `json:"verifier"`
	VerifiedAt  time.Time            `json:"verified_at"`
	Signer      dto.EmployeeResponse `json:"signer"`
	SignedAt    time.Time            `json:"signed_at"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func NewDocumentStatusResponse(document *entity.Document) *DocumentStatusResponse {
	return &DocumentStatusResponse{
		ID:          document.ID,
		Description: document.Description,
		RegisterID:  document.RegisterID,
		Stage:       document.Stage.Status,
		Verifier:    *dto.NewEmployeeResponse(&document.Verifier),
		VerifiedAt:  document.VerifiedAt,
		Signer:      *dto.NewEmployeeResponse(&document.Signer),
		SignedAt:    document.SignedAt,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}
}

type BriefDocumentResponse struct {
	ID          string                `json:"id"`
	Description string                `json:"description"`
	RegisterID  uint                  `json:"register"`
	Applicant   dto.ApplicantResponse `json:"applicant"`
	Stage       string                `json:"stage"`
	Template    string                `json:"template"`
}

func NewBriefDocumentResponse(document *entity.Document) *BriefDocumentResponse {
	return &BriefDocumentResponse{
		ID:          document.ID,
		Description: document.Description,
		RegisterID:  document.RegisterID,
		Applicant:   *dto.NewApplicantResponse(&document.Applicant),
		Stage:       document.Stage.Status,
		Template:    document.Template.Name,
	}
}

type BriefDocumentsResponse []BriefDocumentResponse

func NewBriefDocumentsResponse(documents *entity.Documents) *BriefDocumentsResponse {
	var documentsResponse BriefDocumentsResponse
	for _, document := range *documents {
		documentsResponse = append(documentsResponse, *NewBriefDocumentResponse(&document))
	}

	return &documentsResponse
}

type DocumentUpdateRequest struct {
	RegisterID  uint   `json:"register"`
	Description string `json:"description"`
}

func (d *DocumentUpdateRequest) ToEntity() *entity.Document {
	return &entity.Document{
		RegisterID:  d.RegisterID,
		Description: d.Description,
	}
}

type VerifyDocumentRequest struct {
	RegisterID  uint   `json:"register_id"`
	Description string `json:"description"`
}

func (v *VerifyDocumentRequest) ToEntity() *entity.Document {
	return &entity.Document{
		RegisterID:  v.RegisterID,
		Description: v.Description,
	}
}

type FieldUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (f *FieldUpdateRequest) ToEntity(docID string) *entity.DocumentField {
	return &entity.DocumentField{
		Model: gorm.Model{
			ID: f.ID,
		},
		DocumentID: docID,
		Value:      f.Value,
	}
}

type FieldsUpdateRequest struct {
	Fields []FieldUpdateRequest `json:"fields" validate:"dive"`
}

func (f *FieldsUpdateRequest) ToEntity(docID string) *entity.DocumentFields {
	var fields entity.DocumentFields
	for _, field := range f.Fields {
		fields = append(fields, *field.ToEntity(docID))
	}

	return &fields
}
