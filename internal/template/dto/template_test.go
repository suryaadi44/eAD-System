package dto

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

func TestTemplateRequest_ToEntity(t *testing.T) {
	tests := []struct {
		name string
		tr   TemplateRequest
		want *entity.Template
	}{
		{
			name: "All fields are filled",
			tr: TemplateRequest{
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Keys:         []string{"key1", "key2"},
			},
			want: &entity.Template{
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Fields: []entity.TemplateField{
					{
						Key: "key1",
					},
					{
						Key: "key2",
					},
				},
			},
		},
		{
			name: "Partial fields are filled",
			tr: TemplateRequest{
				Name: "Template 1",
			},
			want: &entity.Template{
				Name: "Template 1",
			},
		},
		{
			name: "All fields are empty",
			tr:   TemplateRequest{},
			want: &entity.Template{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToEntity(); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewTemplateResponse(t *testing.T) {
	type args struct {
		template *entity.Template
	}
	tests := []struct {
		name string
		args args
		want *TemplateResponse
	}{
		{
			name: "All fields are filled",
			args: args{
				template: &entity.Template{
					Model: gorm.Model{
						ID: 1,
					},
					Name:         "Template 1",
					MarginTop:    10,
					MarginBottom: 10,
					MarginLeft:   10,
					MarginRight:  10,
					Fields: []entity.TemplateField{
						{
							Key: "key1",
						},
						{
							Key: "key2",
						},
					},
				},
			},
			want: &TemplateResponse{
				ID:           1,
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Keys: KeysResponse{
					{
						Key: "key1",
					},
					{
						Key: "key2",
					},
				},
			},
		},
		{
			name: "Partial fields are filled",
			args: args{
				template: &entity.Template{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Template 1",
				},
			},
			want: &TemplateResponse{
				ID:   1,
				Name: "Template 1",
			},
		},
		{
			name: "All fields are empty",
			args: args{
				template: &entity.Template{},
			},
			want: &TemplateResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTemplateResponse(tt.args.template); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewTemplatesResponse(t *testing.T) {
	type args struct {
		templates *entity.Templates
	}
	tests := []struct {
		name string
		args args
		want *TemplatesResponse
	}{
		{
			name: "All fields are filled",
			args: args{
				templates: &entity.Templates{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Name:         "Template 1",
						MarginTop:    10,
						MarginBottom: 10,
						MarginLeft:   10,
						MarginRight:  10,
						Fields: []entity.TemplateField{
							{
								Key: "key1",
							},
						},
					},
				},
			},
			want: &TemplatesResponse{
				{
					ID:           1,
					Name:         "Template 1",
					MarginTop:    10,
					MarginBottom: 10,
					MarginLeft:   10,
					MarginRight:  10,
					Keys: KeysResponse{
						{
							Key: "key1",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTemplatesResponse(tt.args.templates); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
