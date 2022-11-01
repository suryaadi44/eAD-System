package dto

import "github.com/suryaadi44/eAD-System/pkg/entity"

type TemplateRequest struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	MarginTop    uint     `json:"margin_top"`
	MarginBottom uint     `json:"margin_bottom"`
	MarginLeft   uint     `json:"margin_left"`
	MarginRight  uint     `json:"margin_right"`
	Keys         []string `json:"keys"`
}

func (t TemplateRequest) ToEntity() (entity.Template, []entity.TemplateField) {
	template := entity.Template{
		Name:         t.Name,
		Path:         t.Path,
		MarginTop:    t.MarginTop,
		MarginBottom: t.MarginBottom,
		MarginLeft:   t.MarginLeft,
		MarginRight:  t.MarginRight,
	}

	var fields []entity.TemplateField
	for _, key := range t.Keys {
		fields = append(fields, entity.TemplateField{
			Key: key,
		})
	}

	return template, fields
}
