package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/pkg/utils"
)

type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) AddTemplate(ctx context.Context, template *dto.TemplateRequest, file io.Reader, fileName string) error {
	args := m.Called(ctx, template, file, fileName)
	return args.Error(0)
}

func (m *MockDocumentService) GetAllTemplate(ctx context.Context) (*dto.TemplatesResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*dto.TemplatesResponse), args.Error(1)
}

func (m *MockDocumentService) GetTemplateDetail(ctx context.Context, templateId uint) (*dto.TemplateResponse, error) {
	args := m.Called(ctx, templateId)
	return args.Get(0).(*dto.TemplateResponse), args.Error(1)
}

func (m *MockDocumentService) AddDocument(ctx context.Context, document dto.DocumentRequest, userID string) (string, error) {
	args := m.Called(ctx, document, userID)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
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

func (m *MockDocumentService) VerifyDocument(ctx context.Context, documentID string, verifierID string) error {
	args := m.Called(ctx, documentID, verifierID)
	return args.Error(0)
}

func (m *MockDocumentService) SignDocument(ctx context.Context, documentID string, signerID string) error {
	args := m.Called(ctx, documentID, signerID)
	return args.Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GetClaims(c *echo.Context) jwt.MapClaims {
	args := m.Called(c)
	return args.Get(0).(jwt.MapClaims)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(a0 interface{}) error {
	args := m.Called(a0)
	return args.Error(0)
}

type TestSuiteDocumentController struct {
	suite.Suite
	mockDocumentService *MockDocumentService
	mockJWTService      *MockJWTService
	mockValidator       *MockValidator
	documentController  *DocumentController
	echoApp             *echo.Echo
}

func (s *TestSuiteDocumentController) SetupTest() {
	s.mockDocumentService = new(MockDocumentService)
	s.mockJWTService = new(MockJWTService)
	s.mockValidator = new(MockValidator)
	s.documentController = NewDocumentController(s.mockDocumentService, s.mockJWTService)
	s.echoApp = echo.New()
	s.echoApp.Validator = s.mockValidator
}

func (s *TestSuiteDocumentController) TearDownTest() {
	s.mockDocumentService = nil
	s.mockJWTService = nil
	s.mockValidator = nil
	s.documentController = nil
	s.echoApp = nil
}

func (s *TestSuiteDocumentController) TestAddTemplate() {
	for _, tc := range []struct {
		Name            string
		RequestBody     interface{}
		FunctionError   error
		JWTReturn       jwt.MapClaims
		ValidationError error
		ExpectedStatus  int
		ExpectedBody    echo.Map
		ExpectedError   error
	}{
		{
			Name: "Success",
			RequestBody: dto.TemplateRequest{
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Keys:         []string{"key1", "key2"},
			},
			JWTReturn:      jwt.MapClaims{"role": float64(3)},
			ExpectedStatus: 200,
			ExpectedBody:   echo.Map{"message": "success adding template"},
		},
		{
			Name: "Failed adding template : insufficient role",
			RequestBody: dto.TemplateRequest{
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Keys:         []string{"key1", "key2"},
			},
			JWTReturn:      jwt.MapClaims{"role": float64(1)},
			ExpectedStatus: http.StatusForbidden,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:           "Failed adding template : invalid request body",
			RequestBody:    "invalid request body",
			JWTReturn:      jwt.MapClaims{"role": float64(3)},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name: "Failed adding template : validation error",
			RequestBody: dto.TemplateRequest{
				Name: "Template 1",
			},
			JWTReturn:       jwt.MapClaims{"role": float64(3)},
			ValidationError: echo.NewHTTPError(http.StatusBadRequest, "validation error"),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedError:   errors.New("validation error"),
		},
		{
			Name: "Failed adding template : service error",
			RequestBody: dto.TemplateRequest{
				Name:         "Template 1",
				MarginTop:    10,
				MarginBottom: 10,
				MarginLeft:   10,
				MarginRight:  10,
				Keys:         []string{"key1", "key2"},
			},
			JWTReturn:      jwt.MapClaims{"role": float64(3)},
			FunctionError:  errors.New("service error"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("service error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			var r *http.Request

			// check if request body is instance of dto.TemplateRequest
			if _, ok := tc.RequestBody.(dto.TemplateRequest); ok {
				// create multipart form data request body
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("name", tc.RequestBody.(dto.TemplateRequest).Name)
				writer.WriteField("margin_top", fmt.Sprint(tc.RequestBody.(dto.TemplateRequest).MarginTop))
				writer.WriteField("margin_bottom", fmt.Sprint(tc.RequestBody.(dto.TemplateRequest).MarginBottom))
				writer.WriteField("margin_left", fmt.Sprint(tc.RequestBody.(dto.TemplateRequest).MarginLeft))
				writer.WriteField("margin_right", fmt.Sprint(tc.RequestBody.(dto.TemplateRequest).MarginRight))
				for _, field := range tc.RequestBody.(dto.TemplateRequest).Keys {
					writer.WriteField("keys[]", field)
				}

				// create form-data
				part, err := writer.CreateFormFile("template", "test.html")
				if err != nil {
					s.FailNow("failed to create form file")
				}

				// create file
				file, err := os.Open("../../../template/test.html")
				if err != nil {
					s.FailNow("failed to open test.html")
				}
				defer file.Close()

				// copy file to form-data
				_, err = io.Copy(part, file)
				if err != nil {
					s.FailNow("failed to copy file to form file")
				}

				// close writer
				err = writer.Close()
				if err != nil {
					s.FailNow("failed to close writer")
				}

				r = httptest.NewRequest(http.MethodPost, "/templates", body)
				r.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			} else {
				body, err := json.Marshal(tc.RequestBody)
				s.NoError(err)
				r = httptest.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}

			w := httptest.NewRecorder()

			// create context
			c := s.echoApp.NewContext(r, w)

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("AddTemplate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionError)

			if tc.ValidationError != nil {
				s.mockValidator.On("Validate", mock.Anything).Return(tc.ValidationError)
			} else {
				s.mockValidator.On("Validate", mock.Anything).Return(nil)
			}

			err := s.documentController.AddTemplate(c)

			if tc.ExpectedError != nil {
				s.Equal(echo.NewHTTPError(tc.ExpectedStatus, tc.ExpectedError.Error()), err)
			} else {
				s.NoError(err)

				var response echo.Map
				err := json.Unmarshal(w.Body.Bytes(), &response)
				s.NoError(err)

				s.Equal(tc.ExpectedStatus, w.Result().StatusCode)
				s.Equal(tc.ExpectedBody, response)
			}

			s.TearDownTest()
		})
	}
}

func (s *TestSuiteDocumentController) TestAddTemplate_Form_data_error() {
	s.SetupTest()

	JWTReturn := jwt.MapClaims{"role": float64(3)}
	ExpectedStatus := http.StatusBadRequest
	ExpectedError := utils.ErrBadRequestBody

	r := httptest.NewRequest(http.MethodPost, "/templates", nil)
	r.Header.Set(echo.HeaderContentType, "multipart/form-data")

	w := httptest.NewRecorder()

	c := s.echoApp.NewContext(r, w)

	s.mockJWTService.On("GetClaims", mock.Anything).Return(JWTReturn)
	s.mockDocumentService.On("AddTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.mockValidator.On("Validate", mock.Anything).Return(nil)

	err := s.documentController.AddTemplate(c)

	s.Equal(echo.NewHTTPError(ExpectedStatus, ExpectedError.Error()), err)

	s.TearDownTest()
}

func TestDocumentController(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentController))
}
