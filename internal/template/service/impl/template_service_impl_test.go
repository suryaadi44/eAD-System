package impl

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	dto3 "github.com/suryaadi44/eAD-System/internal/template/dto"
	mockTemplateRepoPkg "github.com/suryaadi44/eAD-System/internal/template/repository/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
	"os"
	"testing"
)

type TestSuiteTemplateService struct {
	suite.Suite
	mockTemplateRepository *mockTemplateRepoPkg.MockTemplateRepository
	templateService        *TemplateServiceImpl
}

func (s *TestSuiteTemplateService) SetupTest() {
	s.mockTemplateRepository = new(mockTemplateRepoPkg.MockTemplateRepository)
	s.templateService = &TemplateServiceImpl{
		templateRepository: s.mockTemplateRepository,
	}
}

func (s *TestSuiteTemplateService) TearDownTest() {
	s.mockTemplateRepository = nil
	s.templateService = nil
}

func (s *TestSuiteTemplateService) TestAddTemplateToRepo_Success() {
	file, err := os.Open("../../../../template/test.html")
	if err != nil {
		s.Fail("Error when opening file")
	}

	defer file.Close()

	tmp := &entity.Template{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "Test Template",
		Path:         "test.html",
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
		Fields: []entity.TemplateField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Key: "field1",
			},
		},
	}

	s.mockTemplateRepository.On("AddTemplate", mock.Anything, mock.Anything).Return(nil)

	err = s.templateService.addTemplateToRepo(context.Background(), tmp)
	s.NoError(err)
}

func (s *TestSuiteTemplateService) TestAddTemplateToRepo_FailRepoError() {
	file, err := os.Open("../../../../template/test.html")
	if err != nil {
		s.Fail("Error when opening file")
	}

	defer file.Close()

	tmp := &entity.Template{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "Test Template",
		Path:         "test.html",
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
		Fields: []entity.TemplateField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Key: "field1",
			},
		},
	}

	s.mockTemplateRepository.On("AddTemplate", mock.Anything, mock.Anything).Return(errors.New("error"))

	err = s.templateService.addTemplateToRepo(context.Background(), tmp)
	s.Error(err)
}

func (s *TestSuiteTemplateService) TestGetAllTemplate_Success() {
	tmp := &entity.Templates{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:         "Test Template",
			Path:         "test.html",
			MarginTop:    10,
			MarginBottom: 10,
			MarginLeft:   10,
			MarginRight:  10,
			Fields: []entity.TemplateField{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Key: "field1",
				},
			},
		},
	}

	expectedReturn := &dto3.TemplatesResponse{
		{
			ID:           1,
			Name:         "Test Template",
			MarginTop:    10,
			MarginBottom: 10,
			MarginLeft:   10,
			MarginRight:  10,
			Keys: dto3.KeysResponse{
				{
					ID:  1,
					Key: "field1",
				},
			},
		},
	}

	s.mockTemplateRepository.On("GetAllTemplate", mock.Anything).Return(tmp, nil)

	actualTmp, err := s.templateService.GetAllTemplate(context.Background())
	s.NoError(err)
	s.Equal(expectedReturn, actualTmp)
}

func (s *TestSuiteTemplateService) TestGetAllTemplate_RepositoryGenericError() {
	s.mockTemplateRepository.On("GetAllTemplate", mock.Anything).Return(&entity.Templates{}, errors.New("error"))

	_, err := s.templateService.GetAllTemplate(context.Background())
	s.Equal(err, errors.New("error"))
}

func (s *TestSuiteTemplateService) TestGetTemplateDetail() {
	tmp := &entity.Template{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "Test Template",
		Path:         "test.html",
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
		Fields: []entity.TemplateField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Key: "field1",
			},
		},
	}

	expectedReturn := &dto3.TemplateResponse{
		ID:           1,
		Name:         "Test Template",
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
		Keys: dto3.KeysResponse{
			{
				ID:  1,
				Key: "field1",
			},
		},
	}

	s.mockTemplateRepository.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(tmp, nil)

	actualTmp, err := s.templateService.GetTemplateDetail(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedReturn, actualTmp)
}

func (s *TestSuiteTemplateService) TestGetTemplateDetail_RepositoryGenericError() {
	s.mockTemplateRepository.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(&entity.Template{}, errors.New("error"))

	_, err := s.templateService.GetTemplateDetail(context.Background(), 1)
	s.Equal(err, errors.New("error"))
}

func TestTemplateService(t *testing.T) {
	suite.Run(t, new(TestSuiteTemplateService))
}
