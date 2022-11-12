package impl

import (
	"bytes"
	"context"
	"errors"
	mockDocumentRepoPkg "github.com/suryaadi44/eAD-System/internal/document/repository/mock"
	error2 "github.com/suryaadi44/eAD-System/pkg/utils"
	"html/template"
	"testing"
	"time"

	dto3 "github.com/suryaadi44/eAD-System/internal/template/dto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	dto2 "github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"gorm.io/gorm"
)

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

type MockTemplateRepository struct {
	mock.Mock
}

func (m *MockTemplateRepository) AddTemplate(ctc context.Context, template *entity.Template) error {
	args := m.Called(ctc, template)
	return args.Error(0)
}

func (m *MockTemplateRepository) GetAllTemplate(ctx context.Context) (*entity.Templates, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.Templates), args.Error(1)
}

func (m *MockTemplateRepository) GetTemplateDetail(ctx context.Context, templateId uint) (*entity.Template, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.Template), args.Error(1)
}

func (m *MockTemplateRepository) GetTemplateFields(ctx context.Context, templateId uint) (*entity.TemplateFields, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*entity.TemplateFields), args.Error(1)
}

type TestSuiteDocumentService struct {
	suite.Suite
	mockDocumentRepository *mockDocumentRepoPkg.MockDocumentRepository
	mockTemplateRepository *MockTemplateRepository
	mockPDFService         *MockPDFService
	mockRenderService      *MockRenderService
	documentService        *DocumentServiceImpl
}

func (s *TestSuiteDocumentService) SetupTest() {
	s.mockDocumentRepository = new(mockDocumentRepoPkg.MockDocumentRepository)
	s.mockTemplateRepository = new(MockTemplateRepository)
	s.mockPDFService = new(MockPDFService)
	s.mockRenderService = new(MockRenderService)
	s.documentService = &DocumentServiceImpl{
		documentRepository: s.mockDocumentRepository,
		templateRepository: s.mockTemplateRepository,
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
	s.NotNil(NewDocumentServiceImpl(s.mockDocumentRepository, s.mockTemplateRepository, s.mockPDFService, s.mockRenderService))
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

	s.mockTemplateRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
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

	s.mockTemplateRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{}, error2.ErrTemplateFieldNotFound)

	id, err := s.documentService.AddDocument(context.Background(), doc, "123")
	s.Equal(err, error2.ErrTemplateFieldNotFound)
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

	s.mockTemplateRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
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
	s.Equal(err, error2.ErrFieldNotMatch)
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

	s.mockTemplateRepository.On("GetTemplateFields", mock.Anything, uint(1)).Return(&entity.TemplateFields{
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
		RegisterID:  123,
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
		RegisterID:  123,
		Description: "",
		Applicant:   dto2.ApplicantResponse{},
		Template: dto3.TemplateResponse{
			ID:   1,
			Name: "Test Template",
		},
		Fields: dto.FieldsResponse{
			{
				ID:    1,
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
			RegisterID:  123,
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
			RegisterID:  123,
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
			RegisterID:  123,
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
			RegisterID:  123,
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
		RegisterID:  123,
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
		RegisterID:  123,
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
		RegisterID:  123,
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
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{}, error2.ErrDocumentNotFound)

	doc, err := s.documentService.GeneratePDFDocument(context.Background(), "1")
	s.Equal(err, error2.ErrDocumentNotFound)
	s.Nil(doc)
}

func (s *TestSuiteDocumentService) TestGeneratePDFDocument_ErrorGenerateHTMLDocument() {
	s.mockDocumentRepository.On("GetDocument", mock.Anything, mock.Anything).Return(&entity.Document{
		ID:          "1",
		RegisterID:  123,
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
		RegisterID:  123,
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
		ID:         "1",
		RegisterID: 123,
		Register: entity.Register{
			Model: gorm.Model{
				ID: 123,
			},
			Description: "Test Register",
		},
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
		"register":   uint(123),
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
		RegisterID:  123,
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
		RegisterID:  123,
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
		RegisterID:  123,
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
		"register":   uint(123),
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
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     1,
		RegisterID:  1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", nil)

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_ErrorGettingStage() {
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return((*entity.Document)(nil), errors.New("error"))

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", (*dto.VerifyDocumentRequest)(nil))

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_ErrorAlreadyVerified() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     2,
		RegisterID:  1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", nil)

	s.Equal(error2.ErrAlreadyVerified, err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_SuccesAutoGenerateDescription() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:    1,
		RegisterID: 1,
		Template: entity.Template{
			Name: "test",
		},
		Applicant: entity.User{
			Name: "user",
		},
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", &dto.VerifyDocumentRequest{})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_SuccesWithDescriptionOnRequest() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:    1,
		RegisterID: 1,
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", &dto.VerifyDocumentRequest{
		Description: "test",
	})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_SuccesWithRegisterOnRequest() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", &dto.VerifyDocumentRequest{
		RegisterID: 1,
	})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_SuccesWithAutoGenerateRegister() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("AddDocumentRegister", mock.Anything, mock.Anything).Return(uint(1), nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", &dto.VerifyDocumentRequest{})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_SuccesWithErrorAutoGenerateRegister() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("AddDocumentRegister", mock.Anything, mock.Anything).Return(uint(0), errors.New("error"))
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", &dto.VerifyDocumentRequest{})

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestVerifyDocument_RepositoryError() {
	returnedBriefDocumentDetail := &entity.Document{
		StageID:     1,
		RegisterID:  1,
		Description: "test",
	}
	s.mockDocumentRepository.On("GetBriefDocument", mock.Anything, "1").Return(returnedBriefDocumentDetail, nil)
	s.mockDocumentRepository.On("VerifyDocument", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.documentService.VerifyDocument(context.Background(), "1", "1", nil)

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

	s.Equal(error2.ErrAlreadySigned, err)
}

func (s *TestSuiteDocumentService) TestSignDocument_ErrorNotVerified() {
	returnedStage := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(error2.ErrNotVerifiedYet, err)
}

func (s *TestSuiteDocumentService) TestSignDocument_RepositoryError() {
	returnedStage := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "1").Return(&returnedStage, nil)
	s.mockDocumentRepository.On("SignDocument", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.documentService.SignDocument(context.Background(), "1", "1")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_SuccesWitUserRole() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("DeleteDocument", mock.Anything, "documentid").Return(nil)

	err := s.documentService.DeleteDocument(context.Background(), "userid", 1, "documentid")

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_SuccesWithAdminRole() {
	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("DeleteDocument", mock.Anything, "documentid").Return(nil)

	err := s.documentService.DeleteDocument(context.Background(), "userid", 2, "documentid")

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_ErrorGettingApplicantID() {
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return((*string)(nil), errors.New("error"))

	err := s.documentService.DeleteDocument(context.Background(), "userid", 1, "documentid")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_ErrorRoleNotSufficentToDeleteOtherUserDocument() {
	userIDReturned := "userid2"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	err := s.documentService.DeleteDocument(context.Background(), "userid", 1, "documentid")

	s.Equal(error2.ErrDidntHavePermission, err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_ErrorGettingDocumentStage() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return((*int)(nil), errors.New("error"))

	err := s.documentService.DeleteDocument(context.Background(), "userid", 1, "documentid")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestDeleteDocument_ErrorDocumentAlreadySigned() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	stageReturned := 3
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	err := s.documentService.DeleteDocument(context.Background(), "userid", 1, "documentid")

	s.Equal(error2.ErrAlreadySigned, err)
}

func (s *TestSuiteDocumentService) TestUpdateDocument_Success() {
	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("UpdateDocument", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.UpdateDocument(context.Background(), &dto.DocumentUpdateRequest{}, "documentid")

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestUpdateDocument_ErrorGettingDocumentStage() {
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return((*int)(nil), errors.New("error"))

	err := s.documentService.UpdateDocument(context.Background(), &dto.DocumentUpdateRequest{}, "documentid")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestUpdateDocument_ErrorDocumentAlreadySigned() {
	stageReturned := 3
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	err := s.documentService.UpdateDocument(context.Background(), &dto.DocumentUpdateRequest{}, "documentid")

	s.Equal(error2.ErrAlreadySigned, err)
}

func (s *TestSuiteDocumentService) TestUpdateDocument_ErrorDocumentAlreadyVerified() {
	stageReturned := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	err := s.documentService.UpdateDocument(context.Background(), &dto.DocumentUpdateRequest{}, "documentid")

	s.Equal(error2.ErrAlreadyVerified, err)
}

func (s *TestSuiteDocumentService) TestUpdateDocument_RepositoryError() {
	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("UpdateDocument", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.documentService.UpdateDocument(context.Background(), &dto.DocumentUpdateRequest{}, "documentid")

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_SuccessWithAdminAccess() {
	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("UpdateDocumentFields", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 3, "documentid", &dto.FieldsUpdateRequest{})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_SuccessWithApplicantAccess() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	stageReturned := 1
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	s.mockDocumentRepository.On("UpdateDocumentFields", mock.Anything, mock.Anything).Return(nil)

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.NoError(err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_ErrorGettingApplicantID() {
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return((*string)(nil), errors.New("error"))

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_ErrorRoleNotSufficentToUpdateOtherUserDocument() {
	userIDReturned := "userid2"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.Equal(error2.ErrDidntHavePermission, err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_ErrorGettingDocumentStage() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return((*int)(nil), errors.New("error"))

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.Equal(errors.New("error"), err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_ErrorDocumentAlreadySigned() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	stageReturned := 3
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.Equal(error2.ErrAlreadySigned, err)
}

func (s *TestSuiteDocumentService) TestUpdateDocumentFields_ErrorDocumentAlreadyVerified() {
	userIDReturned := "userid"
	s.mockDocumentRepository.On("GetApplicantID", mock.Anything, "documentid").Return(&userIDReturned, nil)

	stageReturned := 2
	s.mockDocumentRepository.On("GetDocumentStage", mock.Anything, "documentid").Return(&stageReturned, nil)

	err := s.documentService.UpdateDocumentFields(context.Background(), "userid", 1, "documentid", &dto.FieldsUpdateRequest{})

	s.Equal(error2.ErrAlreadyVerified, err)
}

func TestDocumentService(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentService))
}
