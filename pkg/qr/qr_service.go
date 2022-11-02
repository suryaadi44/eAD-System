package qr

type CodeService interface {
	GenerateQRCode(documentID string) ([]byte, error)
	GenerateBase64QRCode(documentID string) (string, error)
}
