package dto

import "github.com/suryaadi44/eAD-System/pkg/entity"

type TemplateRequest struct {
	Name         string   `form:"name" validate:"required"`
	MarginTop    uint     `form:"margin_top" validate:"required,gte=0"`
	MarginBottom uint     `form:"margin_bottom" validate:"required,gte=0"`
	MarginLeft   uint     `form:"margin_left" validate:"required,gte=0"`
	MarginRight  uint     `form:"margin_right" validate:"required,gte=0"`
	Keys         []string `form:"keys[]" validate:"required"`
}

func (t TemplateRequest) ToEntity() *entity.Template {
	template := entity.Template{
		Name:         t.Name,
		MarginTop:    t.MarginTop,
		MarginBottom: t.MarginBottom,
		MarginLeft:   t.MarginLeft,
		MarginRight:  t.MarginRight,
	}

	var fields entity.TemplateFields
	for _, key := range t.Keys {
		fields = append(fields, entity.TemplateField{
			Key: key,
		})
	}

	template.Fields = fields

	return &template
}

type TemplateResponse struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	MarginTop    uint     `json:"margin_top"`
	MarginBottom uint     `json:"margin_bottom"`
	MarginLeft   uint     `json:"margin_left"`
	MarginRight  uint     `json:"margin_right"`
	Keys         []string `json:"keys"`
}

type TemplatesResponse []TemplateResponse

func NewTemplateResponse(template *entity.Template) *TemplateResponse {
	var keys []string
	for _, field := range template.Fields {
		keys = append(keys, field.Key)
	}

	return &TemplateResponse{
		ID:           template.ID,
		Name:         template.Name,
		MarginTop:    template.MarginTop,
		MarginBottom: template.MarginBottom,
		MarginLeft:   template.MarginLeft,
		MarginRight:  template.MarginRight,
		Keys:         keys,
	}
}

func NewTemplatesResponse(templates *entity.Templates) *TemplatesResponse {
	var responses TemplatesResponse
	for _, template := range *templates {
		responses = append(responses, *NewTemplateResponse(&template))
	}

	return &responses
}
