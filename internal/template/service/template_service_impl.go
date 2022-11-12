package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/suryaadi44/eAD-System/internal/template/dto"
	"github.com/suryaadi44/eAD-System/internal/template/repository"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type TemplateServiceImpl struct {
	templateRepository repository.TemplateRepository
}

func NewTemplateServiceImpl(templateRepository repository.TemplateRepository) TemplateService {
	return &TemplateServiceImpl{
		templateRepository: templateRepository,
	}
}

func (t *TemplateServiceImpl) AddTemplate(ctx context.Context, template *dto.TemplateRequest, file io.Reader, fileName string) error {
	path, err := t.writeTemplateFile(file, fileName)
	if err != nil {
		return err
	}

	templateEntity := template.ToEntity()
	templateEntity.Path = path

	return t.addTemplateToRepo(ctx, templateEntity)
}

func (t *TemplateServiceImpl) addTemplateToRepo(ctx context.Context, template *entity.Template) error {
	return t.templateRepository.AddTemplate(ctx, template)
}

func (*TemplateServiceImpl) writeTemplateFile(file io.Reader, fileName string) (string, error) {
	newFileName := fmt.Sprint(time.Now().UnixNano(), "-", fileName)
	path := filepath.Join("./template", newFileName)

	// check if file already exist
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("file '%s' already exist", newFileName)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	if err = dst.Close(); err != nil {
		return "", err
	}

	return path, nil
}

func (t *TemplateServiceImpl) GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error) {
	templates, err := t.templateRepository.GetAllTemplate(ctx)
	if err != nil {
		return nil, err
	}

	templateResponse := dto.NewTemplatesResponse(templates)

	return templateResponse, nil
}

func (t *TemplateServiceImpl) GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error) {
	tmpl, err := t.templateRepository.GetTemplateDetail(ctx, templateId)
	if err != nil {
		return nil, err
	}

	templateResponse := dto.NewTemplateResponse(tmpl)

	return templateResponse, nil
}
