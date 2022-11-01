package dto

import "github.com/suryaadi44/eAD-System/pkg/entity"

type DocumentRequest struct {
	Register   string        `json:"register" validate:"required"`
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
		Register:   d.Register,
		TemplateID: d.TemplateID,
		Fields:     fields,
	}
}
