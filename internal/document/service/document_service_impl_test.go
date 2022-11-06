package service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) AddTemplate(ctc context.Context, template *entity.Template) error {
	args := m.Called(ctc, template)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetAllTemplate(ctx context.Context) (*entity.Templates, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.Templates), args.Error(1)
}

func (m *MockDocumentRepository) GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.Template), args.Error(1)
}

func (m *MockDocumentRepository) GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.TemplateFields), args.Error(1)
}

func (m *MockDocumentRepository) AddDocument(ctx context.Context, document *entity.Document) (string, error) {
	args := m.Called(ctx, document)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentRepository) GetDocument(ctx context.Context, documentID string) (*entity.Document, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetDocumentStatus(ctx context.Context, documentID string) (*entity.Document, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetApplicantID(ctx context.Context, documentID string) (*string, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockDocumentRepository) GetDocumentStage(ctx context.Context, documentID string) (*int, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDocumentRepository) VerifyDocument(ctx context.Context, document *entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockDocumentRepository) SignDocument(ctx context.Context, document *entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) GeneratePDF(data *bytes.Buffer, marginTop uint, marginBottom uint, marginLeft uint, marginRight uint) ([]byte, error) {
	args := m.Called(data, marginTop, marginBottom, marginLeft, marginRight)
	return args.Get(0).([]byte), args.Error(1)
}

type MockRenderService struct {
	mock.Mock
}

func (m *MockRenderService) GenerateSignature(signer entity.User) (*template.HTML, error) {
	args := m.Called(signer)
	return args.Get(0).(*template.HTML), args.Error(1)
}

func (m *MockRenderService) GenerateFooter(document *entity.Document) (*template.HTML, error) {
	args := m.Called(document)
	return args.Get(0).(*template.HTML), args.Error(1)
}

func (m *MockRenderService) GenerateHTMLDocument(docTemplate *entity.Template, data *map[string]interface{}) (*bytes.Buffer, error) {
	args := m.Called(docTemplate, data)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

type TestSuiteDocumentService struct {
	suite.Suite
	mockDocumentRepository *MockDocumentRepository
	mockPDFService         *MockPDFService
	mockRenderService      *MockRenderService
	documentService        *DocumentServiceImpl
}

func (s *TestSuiteDocumentService) SetupTest() {
	s.mockDocumentRepository = new(MockDocumentRepository)
	s.mockPDFService = new(MockPDFService)
	s.mockRenderService = new(MockRenderService)
	s.documentService = &DocumentServiceImpl{
		documentRepository: s.mockDocumentRepository,
		pdfService:         s.mockPDFService,
		renderService:      s.mockRenderService,
	}
}

func (s *TestSuiteDocumentService) TearDownTest() {
	s.mockDocumentRepository = nil
	s.mockPDFService = nil
	s.mockRenderService = nil
	s.documentService = nil
}

func (s *TestSuiteDocumentService) TestAddTemplateToRepo_Success() {
	file, err := os.Open("../../../template/test.html")
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

	s.mockDocumentRepository.On("AddTemplate", mock.Anything, mock.Anything).Return(nil)

	err = s.documentService.addTemplateToRepo(context.Background(), tmp)
	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestAddTemplateToRepo_FailRepoError() {
	file, err := os.Open("../../../template/test.html")
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

	s.mockDocumentRepository.On("AddTemplate", mock.Anything, mock.Anything).Return(errors.New("error"))

	err = s.documentService.addTemplateToRepo(context.Background(), tmp)
	s.Error(err)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentService))
}
