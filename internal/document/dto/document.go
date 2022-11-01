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
