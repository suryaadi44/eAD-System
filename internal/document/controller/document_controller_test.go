package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	mockDocumentServicePkg "github.com/suryaadi44/eAD-System/internal/document/service/mock"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	mockJwtServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/jwt_service/mock"
	mockValidatorPkg "github.com/suryaadi44/eAD-System/pkg/utils/validation/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	tmpDto "github.com/suryaadi44/eAD-System/internal/template/dto"
	userDto "github.com/suryaadi44/eAD-System/internal/user/dto"
)

type TestSuiteDocumentController struct {
	suite.Suite
	mockDocumentService *mockDocumentServicePkg.MockDocumentService
	mockJWTService      *mockJwtServicePkg.MockJWTService
	mockValidator       *mockValidatorPkg.MockValidator
	documentController  *DocumentController
	echoApp             *echo.Echo
}

func (s *TestSuiteDocumentController) SetupTest() {
	s.mockDocumentService = new(mockDocumentServicePkg.MockDocumentService)
	s.mockJWTService = new(mockJwtServicePkg.MockJWTService)
	s.mockValidator = new(mockValidatorPkg.MockValidator)
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
				TemplateID: 1,
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
				TemplateID: 1,
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
				TemplateID: 1,
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
				TemplateID: 1,
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
				TemplateID: 1,
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
				TemplateID: 1,
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
				RegisterID:  0,
				Description: "",
				Applicant: userDto.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   tmpDto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   userDto.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     userDto.EmployeeResponse{},
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
					"register":    float64(0),
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
				RegisterID:  0,
				Description: "",
				Applicant: userDto.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   tmpDto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   userDto.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     userDto.EmployeeResponse{},
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
				RegisterID:  0,
				Description: "",
				Applicant: userDto.ApplicantResponse{
					ID:       "1",
					Username: "",
					Name:     "",
				},
				Template:   tmpDto.TemplateResponse{},
				Fields:     nil,
				Stage:      "",
				Verifier:   userDto.EmployeeResponse{},
				VerifiedAt: time.Time{},
				Signer:     userDto.EmployeeResponse{},
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
					"register":    float64(0),
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
					RegisterID:  123,
					Applicant: userDto.ApplicantResponse{
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
						"register":    float64(123),
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
					RegisterID:  123,
					Applicant: userDto.ApplicantResponse{
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
						"register":    float64(123),
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
				RegisterID:  123,
				Stage:       "applied",
				Verifier:    userDto.EmployeeResponse{},
				VerifiedAt:  time.Time{},
				Signer:      userDto.EmployeeResponse{},
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
					"register":    float64(123),
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
			s.mockDocumentService.On("GetApplicantID", mock.Anything, "1").Return(&tc.ServiceReturn, tc.ServiceError)
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
			s.mockDocumentService.On("VerifyDocument", mock.Anything, "1", mock.Anything, mock.Anything).Return(tc.ServiceError)

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
			s.mockDocumentService.On("DeleteDocument", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceError)

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

func (s *TestSuiteDocumentController) TestUpdateDocument() {
	for _, tc := range []struct {
		Name                string
		RequestBody         *dto.DocumentUpdateRequest
		RequestContentTypes string
		ServiceError        error
		JWTReturn           jwt.MapClaims
		ExpectedStatus      int
		ExpectedBody        echo.Map
		ExpectedError       error
	}{
		{
			Name: "Success to update document",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success updating document",
			},
			ExpectedError: nil,
		},
		{
			Name: "Failed to update document : role not sufficient to update document",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrDidntHavePermission,
			JWTReturn: jwt.MapClaims{
				"role": float64(1),
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:                "Failed to update document : invalid request body",
			RequestBody:         nil,
			RequestContentTypes: "",
			ServiceError:        utils.ErrBadRequestBody,
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name: "Failed to update document : generic service error",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        errors.New("generic error"),
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name: "Failed to update document : document not found",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrDocumentNotFound,
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name: "Failed to update document : document already signed",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrAlreadySigned,
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrAlreadySigned,
		},
		{
			Name: "Failed to update document : document already Verified",
			RequestBody: &dto.DocumentUpdateRequest{
				RegisterID:  123,
				Description: "description",
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrAlreadyVerified,
			JWTReturn: jwt.MapClaims{
				"role": float64(3),
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrAlreadyVerified,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("PUT", "/documents", bytes.NewReader(jsonBody))
			r.Header.Set("Content-Type", tc.RequestContentTypes)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockDocumentService.On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceError)

			err = s.documentController.UpdateDocument(c)

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

func (s *TestSuiteDocumentController) TestUpdateDocumentFields() {
	for _, tc := range []struct {
		Name                string
		RequestBody         *dto.FieldsUpdateRequest
		RequestContentTypes string
		ServiceError        error
		JWTReturn           jwt.MapClaims
		ValidationErr       error
		ExpectedStatus      int
		ExpectedBody        echo.Map
		ExpectedError       error
	}{
		{
			Name: "Successfully update document fields",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success updating document fields",
			},
			ExpectedError: nil,
		},
		{
			Name:                "Failed to update document fields : bad request body",
			RequestBody:         nil,
			RequestContentTypes: "",
			ServiceError:        nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name: "Failed to update document fields : failed to validate request body",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID: 1,
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        nil,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  echo.NewHTTPError(http.StatusBadRequest, "Value is required"),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("Value is required"),
		},
		{
			Name: "Failed to update document fields : no document found",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrDocumentNotFound,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDocumentNotFound,
		},
		{
			Name: "Failed to update document fields : err already verified",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrAlreadyVerified,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrAlreadyVerified,
		},
		{
			Name: "Failed to update document fields : err already verified",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrAlreadySigned,
			JWTReturn: jwt.MapClaims{
				"role":    float64(3),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrAlreadySigned,
		},
		{
			Name: "Failed to update document fields : err user role not sufficient to update document fields of other user",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        utils.ErrDidntHavePermission,
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name: "Failed to update document fields : generic error from service",
			RequestBody: &dto.FieldsUpdateRequest{
				Fields: []dto.FieldUpdateRequest{
					{
						ID:    1,
						Value: "value1",
					},
				},
			},
			RequestContentTypes: "application/json",
			ServiceError:        errors.New("generic error"),
			JWTReturn: jwt.MapClaims{
				"role":    float64(1),
				"user_id": "1",
			},
			ValidationErr:  nil,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("PUT", "/documents", bytes.NewReader(jsonBody))
			r.Header.Set("Content-Type", tc.RequestContentTypes)
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)
			c.SetParamNames("document_id")
			c.SetParamValues("1")

			s.mockJWTService.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockValidator.On("Validate", mock.Anything).Return(tc.ValidationErr)
			s.mockDocumentService.On("UpdateDocumentFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceError)

			err = s.documentController.UpdateDocumentFields(c)

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
