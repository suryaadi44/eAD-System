package impl

import (
	"bytes"
	"github.com/suryaadi44/eAD-System/pkg/utils/html"
	"github.com/suryaadi44/eAD-System/pkg/utils/qr"
	"html/template"

	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type RenderServiceImpl struct {
	codeService qr.CodeService
}

func NewRenderServiceImpl(codeService qr.CodeService) html.RenderService {
	return &RenderServiceImpl{
		codeService: codeService,
	}
}

func (_ *RenderServiceImpl) GenerateSignature(signer entity.User) (*template.HTML, error) {
	tmpl, err := template.ParseFiles("./template/signature/signature.html")
	if err != nil {
		return nil, err
	}

	m := map[string]string{
		"signerPosition": signer.Position,
		"signerName":     signer.Name,
		"signerNIP":      signer.NIP,
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, m); err != nil {
		return nil, err
	}

	templateHTML := template.HTML(buf.String())

	return &templateHTML, nil
}

func (r *RenderServiceImpl) GenerateFooter(document *entity.Document) (*template.HTML, error) {
	tmpl, err := template.ParseFiles("./template/signature/footer.html")
	if err != nil {
		return nil, err
	}

	qrImage, err := r.codeService.GenerateBase64QRCode(document.ID)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		"qrImage": template.URL(qrImage),
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, m); err != nil {
		return nil, err
	}

	templateHTML := template.HTML(buf.String())

	return &templateHTML, nil
}

func (_ *RenderServiceImpl) GenerateHTMLDocument(docTemplate *entity.Template, data *map[string]interface{}) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles(docTemplate.Path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, *data); err != nil {
		return nil, err
	}

	return buf, nil
}
