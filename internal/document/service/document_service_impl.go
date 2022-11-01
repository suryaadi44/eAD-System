package service

import (
	"context"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/repository"
)

type DocumentServiceImpl struct {
	documentRepository repository.DocumentRepository
}

func NewDocumentServiceImpl(documentRepository repository.DocumentRepository) DocumentService {
	return &DocumentServiceImpl{documentRepository}
}

func (d *DocumentServiceImpl) AddTemplate(ctx context.Context, template dto.TemplateRequest) error {
	return nil
}
