package qr

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

type CodeServiceImpl struct {
	basePath string
}

func NewCodeServiceImpl(basePath string) CodeService {
	return &CodeServiceImpl{basePath: basePath}
}

func (c *CodeServiceImpl) GenerateQRCode(documentID string) ([]byte, error) {
	path := c.basePath + documentID + "/status"

	var qrCode []byte
	qrCode, err := qrcode.Encode(path, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (c *CodeServiceImpl) GenerateBase64QRCode(documentID string) (string, error) {
	qrCode, err := c.GenerateQRCode(documentID)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(qrCode)
	base64StringWithHeader := "data:image/png;base64," + base64String

	return base64StringWithHeader, nil
}
