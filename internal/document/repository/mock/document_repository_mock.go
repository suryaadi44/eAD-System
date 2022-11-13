package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/suryaadi44/eAD-System/pkg/entity"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) AddDocument(ctx context.Context, document *entity.Document) (string, error) {
	args := m.Called(ctx, document)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentRepository) GetDocument(ctx context.Context, documentID string) (*entity.Document, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetBriefDocument(ctx context.Context, documentID string) (*entity.Document, error) {
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

func (m *MockDocumentRepository) GetApplicant(ctx context.Context, documentID string) (*entity.User, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*entity.User), args.Error(1)
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

func (m *MockDocumentRepository) DeleteDocument(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockDocumentRepository) UpdateDocument(ctx context.Context, document *entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockDocumentRepository) UpdateDocumentFields(ctx context.Context, documentFields *entity.DocumentFields) error {
	args := m.Called(ctx, documentFields)
	return args.Error(0)
}

func (m *MockDocumentRepository) AddDocumentRegister(ctx context.Context, register *entity.Register) (uint, error) {
	args := m.Called(ctx, register)
	return args.Get(0).(uint), args.Error(1)
}
