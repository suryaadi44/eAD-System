package pdf

import (
	"bytes"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFServiceImpl struct {
}

func NewPDFService() PDFService {
	return &PDFServiceImpl{}
}

func (*PDFServiceImpl) GeneratePDF(data *bytes.Buffer, marginTop uint, marginBottom uint, marginLeft uint, marginRight uint) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.MarginTop.Set(marginTop)
	pdfg.MarginBottom.Set(marginBottom)
	pdfg.MarginLeft.Set(marginLeft)
	pdfg.MarginRight.Set(marginRight)
	pdfg.Dpi.Set(600)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(data.Bytes()))
	pdfg.AddPage(page)

	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
