package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type DocumentServiceImpl struct {
	documentRepository repository.DocumentRepository
}

func NewDocumentServiceImpl(documentRepository repository.DocumentRepository) DocumentService {
	return &DocumentServiceImpl{documentRepository}
}

func (d *DocumentServiceImpl) AddTemplate(ctx context.Context, template dto.TemplateRequest, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	path := filepath.Join("./template", file.Filename)

	// check if file already exist
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file '%s' already exist", file.Filename)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	if err = dst.Close(); err != nil {
		return err
	}

	templateEntity := template.ToEntity()
	templateEntity.Path = path

	err = d.documentRepository.AddTemplate(ctx, templateEntity)
	if err != nil {
		return err
	}

	return nil
}
