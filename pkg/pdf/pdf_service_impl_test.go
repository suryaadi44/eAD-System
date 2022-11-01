package pdf

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewPDFService(t *testing.T) {
	service := NewPDFService()
	assert.NotNil(t, service)
}

func TestPDFServiceImpl_GeneratePDF(t *testing.T) {
	s := NewPDFService()

	dummyHTML := []byte("<html><body><h1>Hello World</h1></body></html>")
	data := bytes.NewBuffer(dummyHTML)

	result, err := s.GeneratePDF(data, 0, 0, 0, 0)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	//write result byte to file
	err = os.WriteFile("../../tmp/test.pdf", result, 0644)
	assert.Nil(t, err)
}
