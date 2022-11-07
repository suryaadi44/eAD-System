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
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	dto2 "github.com/suryaadi44/eAD-System/internal/user/dto"
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

func (m *MockDocumentService) VerifyDocument(ctx context.Context, documentID string, verifierID string) error {
	args := m.Called(ctx, documentID, verifierID)
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

func (s *TestSuiteDocumentController) TestInitRoute() {
	s.NotPanics(func() {
		s.documentController.InitRoute(s.echoApp.Group("/"), s.echoApp.Group("/"))
	})
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
			Name: "Failed adding template : Duplicate template name",
			RequestBody: dto.TemplateRequest{
				Name: "Template 1",
			},
			JWTReturn:      jwt.MapClaims{"role": float64(3)},
			FunctionError:  utils.ErrDuplicateTemplateName,
			ExpectedStatus: http.StatusConflict,
			ExpectedError:  utils.ErrDuplicateTemplateName,
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

func (s *TestSuiteDocumentController) TestGetAllTemplate() {
	for _, tc := range []struct {
		Name           string
		FunctionError  error
		FunctionReturn *dto.TemplatesResponse
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "success",
			FunctionError: nil,
			FunctionReturn: &dto.TemplatesResponse{
				{
					ID:           1,
					Name:         "test",
					MarginTop:    1,
					MarginBottom: 1,
					MarginLeft:   1,
					MarginRight:  1,
					Keys: dto.KeysResponse{
						{
							ID:  1,
							Key: "test",
						},
					},
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting all template",
				"data": []interface{}{
					map[string]interface{}{
						"id":            float64(1),
						"name":          "test",
						"margin_top":    float64(1),
						"margin_bottom": float64(1),
						"margin_left":   float64(1),
						"margin_right":  float64(1),
						"keys": []interface{}{
							map[string]interface{}{
								"id":  float64(1),
								"key": "test",
							},
						},
					},
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "failed to get all template: No template in database",
			FunctionError:  utils.ErrTemplateNotFound,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrTemplateNotFound,
		},
		{
			Name:           "failed to get all template: generic error from service",
			FunctionError:  errors.New("failed to get all template"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("failed to get all template"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/templates", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)

			s.mockDocumentService.On("GetAllTemplate", mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.documentController.GetAllTemplate(c)

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

func (s *TestSuiteDocumentController) TestGetTemplateDetail() {
	for _, tc := range []struct {
		Name           string
		TemplateID     string
		FunctionError  error
		FunctionReturn *dto.TemplateResponse
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "Success",
			TemplateID:    "1",
			FunctionError: nil,
			FunctionReturn: &dto.TemplateResponse{
				ID:           1,
				Name:         "name",
				MarginTop:    0,
				MarginBottom: 0,
				MarginLeft:   0,
				MarginRight:  0,
				Keys:         nil,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting template detail",
				"data": map[string]interface{}{
					"id":            float64(1),
					"name":          "name",
					"margin_top":    float64(0),
					"margin_bottom": float64(0),
					"margin_left":   float64(0),
					"margin_right":  float64(0),
					"keys":          nil,
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "failed to get template detail: invalid template id",
			TemplateID:     "a",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrInvalidTemplateID,
		},
		{
			Name:           "failed to get template detail: template not found",
			TemplateID:     "1",
			FunctionError:  utils.ErrTemplateNotFound,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrTemplateNotFound,
		},
		{
			Name:           "failed to get template detail: generic error from service",
			TemplateID:     "1",
			FunctionError:  errors.New("failed to get template detail"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("failed to get template detail"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/templates", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("template_id")
			c.SetParamValues(tc.TemplateID)

			s.mockDocumentService.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.documentController.GetTemplateDetail(c)

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
		})
	}
}

func (s *TestSuiteDocumentController) TestAddDocument() {
	for _, tc := range []struct {
		Name               string
		RequestContentType string
		RequestBody        *dto.DocumentRequest
		ValidationErr      error
		FunctionError      error
		FunctionReturn     string
		ExpectedStatus     int
		ExpectedBody       echo.Map
		ExpectedError      error
	}{
		{
			Name:               "Success",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  nil,
			FunctionError:  nil,
			FunctionReturn: "1",
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success adding document",
				"data": map[string]interface{}{
					"id": "1",
				},
			},
		},
		{
			Name:               "failed to add document: invalid request body",
			RequestContentType: "application/xml",
			RequestBody:        nil,
			ValidationErr:      nil,
			FunctionError:      nil,
			ExpectedStatus:     http.StatusBadRequest,
			ExpectedBody:       nil,
			ExpectedError:      utils.ErrBadRequestBody,
		},
		{
			Name:               "failed to add document: validation error",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  echo.NewHTTPError(http.StatusBadRequest, "validation error"),
			FunctionError:  nil,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("validation error"),
		},
		{
			Name:               "failed to add document: template not found",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  nil,
			FunctionError:  utils.ErrTemplateNotFound,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrTemplateNotFound,
		},
		{
			Name:               "failed to add document: field not match",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  nil,
			FunctionError:  utils.ErrFieldNotMatch,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrFieldNotMatch,
		},
		{
			Name:               "failed to add document: duplicate register",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  nil,
			FunctionError:  utils.ErrDuplicateRegister,
			ExpectedStatus: http.StatusConflict,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDuplicateRegister,
		},
		{
			Name:               "failed to add document: generic service error",
			RequestContentType: "application/json",
			RequestBody: &dto.DocumentRequest{
				Register:    "123",
				Description: "description",
				TemplateID:  1,
				Fields: dto.FieldsRequest{
					{
						FieldID: 1,
						Value:   "value1",
					},
				},
			},
			ValidationErr:  nil,
			FunctionError:  errors.New("error"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", tc.RequestContentType)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)

			s.mockDocumentService.On("AddDocument", mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)
			s.mockValidator.On("Validate", mock.Anything).Return(tc.ValidationErr)
			s.mockJWTService.On("GetClaims", mock.Anything).Return(jwt.MapClaims{
				"user_id": "1",
			})

			err = s.documentController.AddDocument(c)

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
		})
	}
}

func (s *TestSuiteDocumentController) TestGetDocument() {
	for _, tc := range []struct {
		Name           string
		ID             string
		FunctionError  error
		FunctionReturn *dto.DocumentResponse
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "success get document",
			ID:            "1",
			FunctionError: nil,
			FunctionReturn: &dto.DocumentResponse{
				ID:          "",
				Register:    "",
				Description: "",
				Applicant: dto2.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   dto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   dto2.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     dto2.EmployeeResponse{},
				SignedAt:   time.Time{},
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			},

			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(1),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting document",
				"data": map[string]interface{}{
					"id":          "",
					"register":    "",
					"description": "",
					"applicant": map[string]interface{}{
						"id":       "1",
						"username": "",
						"name":     "",
					},
					"template": map[string]interface{}{
						"id":            float64(0),
						"name":          "",
						"margin_top":    float64(0),
						"margin_bottom": float64(0),
						"margin_left":   float64(0),
						"margin_right":  float64(0),
						"keys":          interface{}(nil),
					},
					"fields":      interface{}(nil),
					"stage":       "",
					"verifier":    map[string]interface{}{},
					"verified_at": "0001-01-01T00:00:00Z",
					"signer":      map[string]interface{}{},
					"signed_at":   "0001-01-01T00:00:00Z",
					"created_at":  "0001-01-01T00:00:00Z",
					"updated_at":  "0001-01-01T00:00:00Z",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "failed to get document: error document not found",
			ID:             "1",
			FunctionError:  utils.ErrDocumentNotFound,
			FunctionReturn: nil,
			JWTReturn:      nil,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:           "failed to get document: generic error from service",
			ID:             "1",
			FunctionError:  errors.New("error"),
			FunctionReturn: nil,
			JWTReturn:      nil,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("error"),
		},
		{
			Name:          "failed to get document: role note sufficient to get other user document",
			ID:            "1",
			FunctionError: nil,
			FunctionReturn: &dto.DocumentResponse{
				ID:          "",
				Register:    "",
				Description: "",
				Applicant: dto2.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   dto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   dto2.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     dto2.EmployeeResponse{},
				SignedAt:   time.Time{},
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": "2",
				"role":    float64(1),
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:          "success to get document: role sufficient to get other user document",
			ID:            "1",
			FunctionError: nil,
			FunctionReturn: &dto.DocumentResponse{
				ID:          "",
				Register:    "",
				Description: "",
				Applicant: dto2.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   dto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   dto2.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     dto2.EmployeeResponse{},
				SignedAt:   time.Time{},
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": "2",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting document",
				"data": map[string]interface{}{
					"id":          "",
					"register":    "",
					"description": "",
					"applicant": map[string]interface{}{
						"id":       "1",
						"username": "",
						"name":     "",
					},
					"template": map[string]interface{}{
						"id":            float64(0),
						"name":          "",
						"margin_top":    float64(0),
						"margin_bottom": float64(0),
						"margin_left":   float64(0),
						"margin_right":  float64(0),
						"keys":          interface{}(nil),
					},
					"fields":      interface{}(nil),
					"stage":       "",
					"verifier":    map[string]interface{}{},
					"verified_at": "0001-01-01T00:00:00Z",
					"signer":      map[string]interface{}{},
					"signed_at":   "0001-01-01T00:00:00Z",
					"created_at":  "0001-01-01T00:00:00Z",
					"updated_at":  "0001-01-01T00:00:00Z",
				},
			},
			ExpectedError: nil,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues(tc.ID)

			s.mockDocumentService.On("GetDocument", mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)
			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)

			err := s.documentController.GetDocument(c)

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

func (s *TestSuiteDocumentController) TestGetBriefDocument() {
	for _, tc := range []struct {
		Name           string
		Page           string
		Limit          string
		FunctionError  error
		FunctionReturn *dto.BriefDocumentsResponse
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "success to get brief document",
			Page:          "1",
			Limit:         "10",
			FunctionError: nil,
			FunctionReturn: &dto.BriefDocumentsResponse{
				{
					ID:          "1",
					Description: "description",
					Register:    "123",
					Applicant: dto2.ApplicantResponse{
						ID:       "1",
						Username: "Username",
						Name:     "name",
					},
					Stage:    "approved",
					Template: "template",
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting document",
				"data": []interface{}{
					map[string]interface{}{
						"id":          "1",
						"description": "description",
						"register":    "123",
						"applicant": map[string]interface{}{
							"id":       "1",
							"username": "Username",
							"name":     "name",
						},
						"stage":    "approved",
						"template": "template",
					},
				},
				"meta": map[string]interface{}{
					"page":  float64(1),
					"limit": float64(10),
				},
			},
			ExpectedError: nil,
		},
		{
			Name:          "Success to get brief document: empty parameter",
			Page:          "",
			Limit:         "",
			FunctionError: nil,
			FunctionReturn: &dto.BriefDocumentsResponse{
				{
					ID:          "1",
					Description: "description",
					Register:    "123",
					Applicant: dto2.ApplicantResponse{
						ID:       "1",
						Username: "Username",
						Name:     "name",
					},
					Stage:    "approved",
					Template: "template",
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting document",
				"data": []interface{}{
					map[string]interface{}{
						"id":          "1",
						"description": "description",
						"register":    "123",
						"applicant": map[string]interface{}{
							"id":       "1",
							"username": "Username",
							"name":     "name",
						},
						"stage":    "approved",
						"template": "template",
					},
				},
				"meta": map[string]interface{}{
					"page":  float64(1),
					"limit": float64(20),
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "failed to get brief document: invalid page",
			Page:           "a",
			Limit:          "",
			FunctionError:  nil,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrInvalidNumber,
		},
		{
			Name:           "failed to get brief document: invalid limit",
			Page:           "",
			Limit:          "a",
			FunctionError:  nil,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrInvalidNumber,
		},
		{
			Name:           "failed to get brief document: no document",
			Page:           "",
			Limit:          "",
			FunctionError:  utils.ErrDocumentNotFound,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:           "failed to get brief document: generec service error",
			Page:           "",
			Limit:          "",
			FunctionError:  errors.New("generic service error"),
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
				"role":    float64(2),
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic service error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)

			w := httptest.NewRecorder()

			q := r.URL.Query()
			q.Add("page", tc.Page)
			q.Add("limit", tc.Limit)
			r.URL.RawQuery = q.Encode()

			c := s.echoApp.NewContext(r, w)

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("GetBriefDocuments", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.documentController.GetBriefDocument(c)

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

func (s *TestSuiteDocumentController) TestGetDocumentStatus() {
	for _, tc := range []struct {
		Name           string
		FunctionError  error
		FunctionReturn *dto.DocumentStatusResponse
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "Success to get document status",
			FunctionError: nil,
			FunctionReturn: &dto.DocumentStatusResponse{
				ID:          "1",
				Description: "description",
				Register:    "123",
				Stage:       "applied",
				Verifier:    dto2.EmployeeResponse{},
				VerifiedAt:  time.Time{},
				Signer:      dto2.EmployeeResponse{},
				SignedAt:    time.Time{},
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success getting document status",
				"data": map[string]interface{}{
					"id":          "1",
					"description": "description",
					"register":    "123",
					"stage":       "applied",
					"verifier":    map[string]interface{}{},
					"verified_at": "0001-01-01T00:00:00Z",
					"signer":      map[string]interface{}{},
					"signed_at":   "0001-01-01T00:00:00Z",
					"created_at":  "0001-01-01T00:00:00Z",
					"updated_at":  "0001-01-01T00:00:00Z",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "Failed to get document status : document not found",
			FunctionError:  utils.ErrDocumentNotFound,
			FunctionReturn: nil,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:           "Failed to get document status : generic error from service",
			FunctionError:  errors.New("generic error"),
			FunctionReturn: nil,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockDocumentService.On("GetDocumentStatus", mock.Anything, "1").Return(tc.FunctionReturn, tc.FunctionError)

			err := s.documentController.GetDocumentStatus(c)

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

func (s *TestSuiteDocumentController) TestGetPDFDocument() {
	for _, tc := range []struct {
		Name           string
		ServiceError   error
		ServiceReturn  string
		PDFError       error
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedError  error
	}{
		{
			Name:          "Success to get pdf document",
			ServiceError:  nil,
			ServiceReturn: "1",
			PDFError:      nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
		},
		{
			Name:          "Failed to get pdf document : document not found while role is 1",
			ServiceError:  utils.ErrDocumentNotFound,
			ServiceReturn: "",
			PDFError:      nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:          "Failed to get pdf document : generic service error while role is 1",
			ServiceError:  errors.New("generic error"),
			ServiceReturn: "",
			PDFError:      nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:          "Failed to get pdf document : role not sufficient to get other user pdf document",
			ServiceError:  nil,
			ServiceReturn: "2",
			PDFError:      nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:          "Failed to get pdf document : document not found while role is employee",
			ServiceError:  nil,
			ServiceReturn: "",
			PDFError:      utils.ErrDocumentNotFound,
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:          "Failed to get pdf document : generic service error while role is employee",
			ServiceError:  nil,
			ServiceReturn: "",
			PDFError:      errors.New("generic error"),
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("generic error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("GetApplicantID", mock.Anything, "1").Return((*string)(&tc.ServiceReturn), tc.ServiceError)
			s.mockDocumentService.On("GeneratePDFDocument", mock.Anything, "1").Return([]byte(nil), tc.PDFError)

			err := s.documentController.GetPDFDocument(c)

			if tc.ExpectedError != nil {
				s.Equal(echo.NewHTTPError(tc.ExpectedStatus, tc.ExpectedError.Error()), err)
			} else {
				s.NoError(err)
			}

			s.TearDownTest()
		})
	}
}

func (s *TestSuiteDocumentController) TestVerifyDocument() {
	for _, tc := range []struct {
		Name           string
		ServiceError   error
		ServiceReturn  string
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "Success to verify document",
			ServiceError:  nil,
			ServiceReturn: "1",
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success verifying document",
			},
			ExpectedError: nil,
		},
		{
			Name:          "Failed to verify document : document not found",
			ServiceError:  utils.ErrDocumentNotFound,
			ServiceReturn: "",
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:          "Failed to verify document : generic service error",
			ServiceError:  errors.New("generic error"),
			ServiceReturn: "",
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:          "Failed to verify document : role not sufficient to verify document",
			ServiceError:  nil,
			ServiceReturn: "2",
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("VerifyDocument", mock.Anything, "1", mock.Anything).Return(tc.ServiceError)

			err := s.documentController.VerifyDocument(c)

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

func (s *TestSuiteDocumentController) TestSignDocument() {
	for _, tc := range []struct {
		Name           string
		ServiceError   error
		ServiceReturn  string
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "Success to sign document",
			ServiceError:  nil,
			ServiceReturn: "1",
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success signing document",
			},
			ExpectedError: nil,
		},
		{
			Name:          "Failed to sign document : document not found",
			ServiceError:  utils.ErrDocumentNotFound,
			ServiceReturn: "",
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:          "Failed to sign document : generic service error",
			ServiceError:  errors.New("generic error"),
			ServiceReturn: "",
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:          "Failed to sign document : role not sufficient to verify document",
			ServiceError:  nil,
			ServiceReturn: "2",
			JWTReturn: jwt.MapClaims{
				"role":    float64(2),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:          "Failed to sign document : document not verified yet",
			ServiceError:  utils.ErrNotVerifiedYet,
			ServiceReturn: "",
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrNotVerifiedYet,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("SignDocument", mock.Anything, "1", mock.Anything).Return(tc.ServiceError)

			err := s.documentController.SignDocument(c)

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

func (s *TestSuiteDocumentController) TestDeleteDocument() {
	for _, tc := range []struct {
		Name           string
		ServiceError   error
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:         "Success to delete document",
			ServiceError: nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success deleting document",
			},
			ExpectedError: nil,
		},
		{
			Name:         "Failed to delete document : document already signed",
			ServiceError: utils.ErrAlreadySigned,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrAlreadySigned,
		},
		{
			Name:         "Failed to delete document : document not found",
			ServiceError: utils.ErrDocumentNotFound,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name:         "Failed to delete document : generic service error",
			ServiceError: errors.New("generic error"),
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:         "Failed to delete document : role not sufficient to delete other user document",
			ServiceError: utils.ErrDidntHavePermission,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("DELETE", "/documents", nil)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("DeleteDocument", mock.Anything, "1", mock.Anything).Return(tc.ServiceError)

			err := s.documentController.DeleteDocument(c)

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

func TestDocumentController(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentController))
}
