package service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	dto2 "github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
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

func (m *MockDocumentRepository) GetBriefDocuments(ctx context.Context, limit int, offset int) (*entity.Documents, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).(*entity.Documents), args.Error(1)
}

func (m *MockDocumentRepository) GetBriefDocumentsByApplicant(ctx context.Context, applicantID string, limit int, offset int) (*entity.Documents, error) {
	args := m.Called(ctx, applicantID, limit, offset)
	return args.Get(0).(*entity.Documents), args.Error(1)
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

func (s *TestSuiteDocumentService) TestNewDocumentServiceImpl() {
	s.NotNil(NewDocumentServiceImpl(s.mockDocumentRepository, s.mockPDFService, s.mockRenderService))
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

func (s *TestSuiteDocumentService) TestGetAllTemplate_Success() {
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

	expectedReturn := &dto.TemplatesResponse{
		{
			ID:           1,
			Name:         "Test Template",
			MarginTop:    10,
			MarginBottom: 10,
			MarginLeft:   10,
			MarginRight:  10,
			Keys: dto.KeysResponse{
				{
					ID:  1,
					Key: "field1",
				},
			},
		},
	}

	s.mockDocumentRepository.On("GetAllTemplate", mock.Anything).Return(tmp, nil)

	actualTmp, err := s.documentService.GetAllTemplate(context.Background())
	s.NoError(err)
	s.Equal(expectedReturn, actualTmp)
}

func (s *TestSuiteDocumentService) TestGetAllTemplate_RepositoryGenericError() {
	s.mockDocumentRepository.On("GetAllTemplate", mock.Anything).Return(&entity.Templates{}, errors.New("error"))

	_, err := s.documentService.GetAllTemplate(context.Background())
	s.Equal(err, errors.New("error"))
}

func (s *TestSuiteDocumentService) TestGetTemplateDetail() {
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

	expectedReturn := &dto.TemplateResponse{
		ID:           1,
		Name:         "Test Template",
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
		Keys: dto.KeysResponse{
			{
				ID:  1,
				Key: "field1",
			},
		},
	}

	s.mockDocumentRepository.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(tmp, nil)

	actualTmp, err := s.documentService.GetTemplateDetail(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedReturn, actualTmp)
}

func (s *TestSuiteDocumentService) TestGetTemplateDetail_RepositoryGenericError() {
	s.mockDocumentRepository.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(&entity.Template{}, errors.New("error"))

	_, err := s.documentService.GetTemplateDetail(context.Background(), 1)
	s.Equal(err, errors.New("error"))
}

func (s *TestSuiteDocumentService) TestAddDocument_Success() {
	doc := &dto.DocumentRequest{
		TemplateID: 1,
		Fields: dto.FieldsRequest{
			{
				FieldID: 1,
				Value:   "value1",
			},
		},
	}

	s.mockDocumentRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
		{
			Model: gorm.Model{
				ID: 1,
			},
			TemplateID: 1,
			Key:        "field1",
		},
	}, nil)
	s.mockDocumentRepository.On("AddDocument", mock.Anything, mock.Anything).Return("123", nil)

	id, err := s.documentService.AddDocument(context.Background(), doc, "123")
	s.NoError(err)
	s.Equal(id, "123")
}

func (s *TestSuiteDocumentService) TestAddDocument_ErrorNoTemplate() {
	doc := &dto.DocumentRequest{
		TemplateID: 1,
		Fields: dto.FieldsRequest{
			{
				FieldID: 1,
				Value:   "value1",
			},
		},
	}

	s.mockDocumentRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{}, utils.ErrTemplateFieldNotFound)

	id, err := s.documentService.AddDocument(context.Background(), doc, "123")
	s.Equal(err, utils.ErrTemplateFieldNotFound)
	s.Equal(id, "")
}

func (s *TestSuiteDocumentService) TestAddDocument_ErrorFieldMissing() {
	doc := &dto.DocumentRequest{
		TemplateID: 1,
		Fields: dto.FieldsRequest{
			{
				FieldID: 1,
				Value:   "value1",
			},
		},
	}

	s.mockDocumentRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
		{
			Model: gorm.Model{
				ID: 1,
			},
			TemplateID: 1,
			Key:        "field1",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			TemplateID: 1,
			Key:        "field2",
		},
	}, nil)

	id, err := s.documentService.AddDocument(context.Background(), doc, "123")
	s.Equal(err, utils.ErrFieldNotMatch)
	s.Equal(id, "")
}

func (s *TestSuiteDocumentService) TestAddDocument_ErrorRepository() {
	doc := &dto.DocumentRequest{
		TemplateID: 1,
		Fields: dto.FieldsRequest{
			{
				FieldID: 1,
				Value:   "value1",
			},
		},
	}

	s.mockDocumentRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
		{
			Model: gorm.Model{
				ID: 1,
			},
			TemplateID: 1,
			Key:        "field1",
		},
	}, nil)
	s.mockDocumentRepository.On("AddDocument", mock.Anything, mock.Anything).Return("", errors.New("error"))

	id, err := s.documentService.AddDocument(context.Background(), doc, "123")
	s.Equal(err, errors.New("error"))
	s.Equal(id, "")
}

func (s *TestSuiteDocumentService) TestGetDocument_Success() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{
		ID:          "1",
		Register:    "",
		Description: "",
		ApplicantID: "",
		Applicant:   entity.User{},
		TemplateID:  1,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Template",
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}, nil)

	expectedReturn := &dto.DocumentResponse{
		ID:          "1",
		Register:    "",
		Description: "",
		Applicant:   dto2.ApplicantResponse{},
		Template: dto.TemplateResponse{
			ID:   1,
			Name: "Test Template",
		},
		Fields: dto.FieldsResponse{
			{
				Key:   "field1",
				Value: "value1",
			},
		},
		Stage:      "",
		Verifier:   dto2.EmployeeResponse{},
		VerifiedAt: time.Time{},
		Signer:     dto2.EmployeeResponse{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	doc, err := s.documentService.GetDocument(context.Background(), "1")
	s.NoError(err)
	s.Equal(expectedReturn, doc)
}

func (s *TestSuiteDocumentService) TestGetDocument_ErrorRepository() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{}, errors.New("error"))

	doc, err := s.documentService.GetDocument(context.Background(), "1")
	s.Equal(err, errors.New("error"))
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestGetBriefDocuments_Success() {
	s.mockDocumentRepository.On("GetBriefDocuments", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Documents{
		{
			ID:          "1",
			Register:    "register",
			Description: "description",
			Applicant: entity.User{
				ID:       "1",
				Username: "username",
				Name:     "name",
			},
			Template: entity.Template{
				Name: "Test Template",
			},
			Stage: entity.Stage{
				ID:     1,
				Status: "Test Stage",
			},
			CreatedAt: time.Time{},
		},
	}, nil)

	expectedReturn := &dto.BriefDocumentsResponse{
		{
			ID:          "1",
			Register:    "register",
			Description: "description",
			Applicant: dto2.ApplicantResponse{
				ID:       "1",
				Username: "username",
				Name:     "name",
			},
			Template: "Test Template",
			Stage:    "Test Stage",
		},
	}

	docs, err := s.documentService.GetBriefDocuments(context.Background(), "1", 3, 0, 0)
	s.NoError(err)
	s.Equal(expectedReturn, docs)
}

func (s *TestSuiteDocumentService) TestGetBriefDocuments_SuccessWithUserRole() {
	s.mockDocumentRepository.On("GetBriefDocumentsByApplicant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&entity.Documents{
		{
			ID:          "1",
			Register:    "register",
			Description: "description",
			Applicant: entity.User{
				ID:       "1",
				Username: "username",
				Name:     "name",
			},
			Template: entity.Template{
				Name: "Test Template",
			},
			Stage: entity.Stage{
				ID:     1,
				Status: "Test Stage",
			},
			CreatedAt: time.Time{},
		},
	}, nil)

	expectedReturn := &dto.BriefDocumentsResponse{
		{
			ID:          "1",
			Register:    "register",
			Description: "description",
			Applicant: dto2.ApplicantResponse{
				ID:       "1",
				Username: "username",
				Name:     "name",
			},
			Template: "Test Template",
			Stage:    "Test Stage",
		},
	}

	docs, err := s.documentService.GetBriefDocuments(context.Background(), "1", 1, 0, 0)
	s.NoError(err)
	s.Equal(expectedReturn, docs)
}

func (s *TestSuiteDocumentService) TestGetBriefDocuments_ErrorRepository() {
	s.mockDocumentRepository.On("GetBriefDocuments", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Documents{}, errors.New("error"))

	docs, err := s.documentService.GetBriefDocuments(context.Background(), "1", 3, 0, 0)
	s.Equal(err, errors.New("error"))
	s.Nil(docs)
}

func (s *TestSuiteDocumentService) TestGetDocumentStatus_Success() {
	s.mockDocumentRepository.On("GetDocumentStatus", mock.Anything, mock.Anything).Return(&entity.Document{
		ID:          "1",
		Register:    "",
		Description: "",
		ApplicantID: "",
		Applicant:   entity.User{},
		TemplateID:  1,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Template",
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}, nil)

	expectedReturn := &dto.DocumentStatusResponse{
		ID:          "1",
		Description: "",
		Register:    "",
		Stage:       "",
		Verifier:    dto2.EmployeeResponse{},
		VerifiedAt:  time.Time{},
		Signer:      dto2.EmployeeResponse{},
		SignedAt:    time.Time{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	doc, err := s.documentService.GetDocumentStatus(context.Background(), "1")
	s.NoError(err)
	s.Equal(expectedReturn, doc)
}

func (s *TestSuiteDocumentService) TestGetDocumentStatus_ErrorRepository() {
	s.mockDocumentRepository.On("GetDocumentStatus", mock.Anything, mock.Anything).Return(&entity.Document{}, errors.New("error"))

	doc, err := s.documentService.GetDocumentStatus(context.Background(), "1")
	s.Equal(err, errors.New("error"))
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestGeneratePDFDocument_Success() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{
		ID:          "1",
		Register:    "",
		Description: "",
		ApplicantID: "",
		Applicant:   entity.User{},
		TemplateID:  1,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Template",
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}, nil)

	html := `<!DOCTYPE html>`
	buf := bytes.NewBufferString(html)

	s.mockRenderService.On("GenerateHTMLDocument", mock.Anything, mock.Anything).Return(buf, nil)
	s.mockPDFService.On("GeneratePDF", buf, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte("pdf"), nil)

	doc, err := s.documentService.GeneratePDFDocument(context.Background(), "1")
	s.NoError(err)
	s.Equal([]byte("pdf"), doc)
}

func (s *TestSuiteDocumentService) TestGeneratePDFDocument_ErrorDocumentNotFound() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{}, utils.ErrDocumentNotFound)

	doc, err := s.documentService.GeneratePDFDocument(context.Background(), "1")
	s.Equal(err, utils.ErrDocumentNotFound)
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestGeneratePDFDocument_ErrorGenerateHTMLDocument() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{
		ID:          "1",
		Register:    "",
		Description: "",
		ApplicantID: "",
		Applicant:   entity.User{},
		TemplateID:  1,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Template",
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}, nil)

	s.mockRenderService.On("GenerateHTMLDocument", mock.Anything, mock.Anything).Return(&bytes.Buffer{}, errors.New("error"))

	doc, err := s.documentService.GeneratePDFDocument(context.Background(), "1")
	s.Equal(errors.New("error"), err)
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestGeneratePDFDocument_ErrorGeneratePDF() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, "1").Return(&entity.Document{
		ID:          "1",
		Register:    "",
		Description: "",
		ApplicantID: "",
		Applicant:   entity.User{},
		TemplateID:  1,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Template",
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}, nil)

	html := `<!DOCTYPE html>`
	buf := bytes.NewBufferString(html)

	s.mockRenderService.On("GenerateHTMLDocument", mock.Anything, mock.Anything).Return(buf, nil)
	s.mockPDFService.On("GeneratePDF", buf, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte(nil), errors.New("error"))

	doc, err := s.documentService.GeneratePDFDocument(context.Background(), "1")
	s.Equal(errors.New("error"), err)
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestFillMapFields_NotSignedYet() {
	doc := &entity.Document{
		ID:          "1",
		Register:    "register",
		Description: "",
		ApplicantID: "",
		Applicant: entity.User{
			ID:       "1",
			Username: "",
			Name:     "",
		},
		TemplateID: 0,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    0,
		Stage:      entity.Stage{},
		VerifierID: "",
		Verifier:   entity.User{},
		VerifiedAt: time.Time{},
		SignerID:   "",
		Signer:     entity.User{},
		SignedAt:   time.Time{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  gorm.DeletedAt{},
	}

	expectedMap := &map[string]interface{}{
		"field1":     "value1",
		"register":   "register",
		"signedDate": "",
		"signature":  "",
		"footer":     "",
	}

	m, err := s.documentService.fillMapFields(doc)

	s.Equal(expectedMap, m)
	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestFillMapFields_SignatureError() {
	now := time.Now()

	doc := &entity.Document{
		ID:          "1",
		Register:    "register",
		Description: "",
		ApplicantID: "",
		Applicant: entity.User{
			ID:       "1",
			Username: "",
			Name:     "",
		},
		TemplateID: 0,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    2,
		Stage:      entity.Stage{},
		VerifierID: "1",
		Verifier: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		VerifiedAt: now,
		SignerID:   "1",
		Signer: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		SignedAt:  now,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	s.mockRenderService.On("GenerateSignature", mock.Anything, mock.Anything).Return((*template.HTML)(nil), errors.New("error"))

	m, err := s.documentService.fillMapFields(doc)
	s.Nil(m)
	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestFillMapFields_FooterError() {
	now := time.Now()

	doc := &entity.Document{
		ID:          "1",
		Register:    "register",
		Description: "",
		ApplicantID: "",
		Applicant: entity.User{
			ID:       "1",
			Username: "",
			Name:     "",
		},
		TemplateID: 0,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    2,
		Stage:      entity.Stage{},
		VerifierID: "1",
		Verifier: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		VerifiedAt: now,
		SignerID:   "1",
		Signer: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		SignedAt:  now,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	templateHtml := template.HTML(`<!DOCTYPE html>`)
	s.mockRenderService.On("GenerateSignature", mock.Anything, mock.Anything).Return(&templateHtml, nil)
	s.mockRenderService.On("GenerateFooter", mock.Anything, mock.Anything).Return((*template.HTML)(nil), errors.New("error"))

	m, err := s.documentService.fillMapFields(doc)
	s.Nil(m)
	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestFillMapFields_Success() {
	now := time.Now()

	doc := &entity.Document{
		ID:          "1",
		Register:    "register",
		Description: "",
		ApplicantID: "",
		Applicant: entity.User{
			ID:       "1",
			Username: "",
			Name:     "",
		},
		TemplateID: 0,
		Template: entity.Template{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Fields: []entity.DocumentField{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DocumentID:      "1",
				TemplateFieldID: 1,
				TemplateField: entity.TemplateField{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "field1",
				},
				Value: "value1",
			},
		},
		StageID:    2,
		Stage:      entity.Stage{},
		VerifierID: "1",
		Verifier: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		VerifiedAt: now,
		SignerID:   "1",
		Signer: entity.User{
			ID:       "1",
			NIP:      "1234567890",
			Username: "",
			Name:     "",
			Position: "position",
		},
		SignedAt:  now,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	templateHtml := template.HTML(`<!DOCTYPE html>`)
	expectedMap := &map[string]interface{}{
		"field1":     "value1",
		"register":   "register",
		"signedDate": now.Format("02 January 2006"),
		"signature":  &templateHtml,
		"footer":     &templateHtml,
	}

	s.mockRenderService.On("GenerateSignature", mock.Anything, mock.Anything).Return(&templateHtml, nil)
	s.mockRenderService.On("GenerateFooter", mock.Anything, mock.Anything).Return(&templateHtml, nil)

	m, err := s.documentService.fillMapFields(doc)

	s.Equal(expectedMap, m)
	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestGetApplicantID_Success() {
	returnedID := "1"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "1").Return(&returnedID, nil)

	id, err := s.documentService.GetApplicantID(context.Background(), "1")

	s.Equal(&returnedID, id)
	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestGetApplicantID_Error() {
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "1").Return((*string)(nil), errors.New("error"))

	id, err := s.documentService.GetApplicantID(context.Background(), "1")

	s.Equal((*string)(nil), id)
	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_Success() {
	returnedStage := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1")

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_ErrorGettingStage() {
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return((*int)(nil), errors.New("error"))

	err := s.documentService.VerifyDocument(context.Background(), "1", "1")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_ErrorAlreadyVerified() {
	returnedStage := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1")

	s.Equal(utils.ErrAlreadyVerified, err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_RepositoryError() {
	returnedStage := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.documentService.VerifyDocument(context.Background(), "1", "1")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestSignDocument_Success() {
	returnedStage := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)
	s.mockDocumentRepository.On("SignDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestSignDocument_ErrorGettingStage() {
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return((*int)(nil), errors.New("error"))

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestSignDocument_ErrorAlreadySigned() {
	returnedStage := 3
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(utils.ErrAlreadySigned, err)
}

func (s *TestSuiteDocumentService) TestSignDocument_ErrorNotVerified() {
	returnedStage := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(utils.ErrNotVerifiedYet, err)
}

func (s *TestSuiteDocumentService) TestSignDocument_RepositoryError() {
	returnedStage := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)
	s.mockDocumentRepository.On("SignDocument", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(errors.New("error"), err)
}

func TestDocumentService(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentService))
}
