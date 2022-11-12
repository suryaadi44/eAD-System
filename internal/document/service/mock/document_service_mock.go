package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
)

type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) AddDocument(ctx context.Context, document *dto.DocumentRequest, userID string) (string, error) {
	args := m.Called(ctx, document, userID)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

func (m *MockDocumentService) GetBriefDocuments(ctx context.Context, applicantID string, role int, page int, limit int) (*dto.BriefDocumentsResponse, error) {
	args := m.Called(ctx, applicantID, role, page, limit)
	return args.Get(0).(*dto.BriefDocumentsResponse), args.Error(1)
}

func (m *MockDocumentService) GetDocumentStatus(ctx context.Context, documentID string) (*dto.DocumentStatusResponse, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*dto.DocumentStatusResponse), args.Error(1)
}

func (m *MockDocumentService) GeneratePDFDocument(ctx context.Context, documentID string) ([]byte, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDocumentService) GetApplicantID(ctx context.Context, documentID string) (*string, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockDocumentService) VerifyDocument(ctx context.Context, documentID string, verifierID string, verifyRequest *dto.VerifyDocumentRequest) error {
	args := m.Called(ctx, documentID, verifierID, verifyRequest)
	return args.Error(0)
}

func (m *MockDocumentService) SignDocument(ctx context.Context, documentID string, signerID string) error {
	args := m.Called(ctx, documentID, signerID)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, userID string, role int, documentID string) error {
	args := m.Called(ctx, userID, role, documentID)
	return args.Error(0)
}

func (m *MockDocumentService) UpdateDocument(ctx context.Context, document *dto.DocumentUpdateRequest, documentID string) error {
	args := m.Called(ctx, document, documentID)
	return args.Error(0)
}

func (m *MockDocumentService) UpdateDocumentFields(ctx context.Context, userID string, role int, documentID string, fields *dto.FieldsUpdateRequest) error {
	args := m.Called(ctx, userID, role, documentID, fields)
	return args.Error(0)
}
