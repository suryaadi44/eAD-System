package pdf

import "bytes"

type PDFService interface {
	GeneratePDF(data *bytes.Buffer, marginTop uint, marginBottom uint, marginLeft uint, marginRight uint) ([]byte, error)
}
