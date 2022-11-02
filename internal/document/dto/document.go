package dto

import (
	"time"

	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type DocumentRequest struct {
	Register    string        `json:"register" validate:"required"`
	Description string        `json:"description" validate:"required"`
	TemplateID  uint          `json:"template_id" validate:"required"`
	Fields      FieldsRequest `json:"fields" validate:"required"`
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
		Register:    d.Register,
		Description: d.Description,
		TemplateID:  d.TemplateID,
		Fields:      fields,
	}
}

type DocumentResponse struct {
	ID          string                `json:"id"`
	Register    string                `json:"register"`
	Description string                `json:"description"`
	Applicant   dto.ApplicantResponse `json:"applicant_id"`
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
		Register:    document.Register,
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
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewFieldResponse(fields *entity.DocumentField) *FieldResponse {
	return &FieldResponse{
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
	Register    string               `json:"register"`
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
		Register:    document.Register,
		Stage:       document.Stage.Status,
		Verifier:    *dto.NewEmployeeResponse(&document.Verifier),
		VerifiedAt:  document.VerifiedAt,
		Signer:      *dto.NewEmployeeResponse(&document.Signer),
		SignedAt:    document.SignedAt,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}
}
