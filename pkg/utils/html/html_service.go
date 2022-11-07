package html

import (
	"bytes"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"html/template"
)

type RenderService interface {
	GenerateSignature(signer entity.User) (*template.HTML, error)
	GenerateFooter(document *entity.Document) (*template.HTML, error)
	GenerateHTMLDocument(docTemplate *entity.Template, data *map[string]interface{}) (*bytes.Buffer, error)
}
