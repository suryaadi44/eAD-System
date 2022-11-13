package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	mockTemplateServicePkg "github.com/suryaadi44/eAD-System/internal/template/service/mock"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	mockJwtServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/jwt_service/mock"
	mockValidatorPkg "github.com/suryaadi44/eAD-System/pkg/utils/validation/mock"
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
	"github.com/suryaadi44/eAD-System/internal/template/dto"
)

type TestSuiteTemplateController struct {
	suite.Suite
	mockTemplateService *mockTemplateServicePkg.MockTemplateService
	mockJWTService      *mockJwtServicePkg.MockJWTService
	mockValidator       *mockValidatorPkg.MockValidator
	templateController  *TemplateController
	echoApp             *echo.Echo
}

func (s *TestSuiteTemplateController) SetupTest() {
	s.mockTemplateService = new(mockTemplateServicePkg.MockTemplateService)
	s.mockJWTService = new(mockJwtServicePkg.MockJWTService)
	s.mockValidator = new(mockValidatorPkg.MockValidator)
	s.templateController = NewTemplateController(s.mockTemplateService, s.mockJWTService)
	s.echoApp = echo.New()
	s.echoApp.Validator = s.mockValidator
}

func (s *TestSuiteTemplateController) TearDownTest() {
	s.mockTemplateService = nil
	s.mockJWTService = nil
	s.mockValidator = nil
	s.templateController = nil
	s.echoApp = nil
}

func (s *TestSuiteTemplateController) TestAddTemplate() {
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
			s.mockTemplateService.On("AddTemplate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionError)

			if tc.ValidationError != nil {
				s.mockValidator.On("Validate", mock.Anything).Return(tc.ValidationError)
			} else {
				s.mockValidator.On("Validate", mock.Anything).Return(nil)
			}

			err := s.templateController.AddTemplate(c)

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

func (s *TestSuiteTemplateController) TestAddTemplate_Form_data_error() {
	s.SetupTest()

	JWTReturn := jwt.MapClaims{"role": float64(3)}
	ExpectedStatus := http.StatusBadRequest
	ExpectedError := utils.ErrBadRequestBody

	r := httptest.NewRequest(http.MethodPost, "/templates", nil)
	r.Header.Set(echo.HeaderContentType, "multipart/form-data")

	w := httptest.NewRecorder()

	c := s.echoApp.NewContext(r, w)

	s.mockJWTService.On("GetClaims", mock.Anything).Return(JWTReturn)
	s.mockTemplateService.On("AddTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.mockValidator.On("Validate", mock.Anything).Return(nil)

	err := s.templateController.AddTemplate(c)

	s.Equal(echo.NewHTTPError(ExpectedStatus, ExpectedError.Error()), err)

	s.TearDownTest()
}

func (s *TestSuiteTemplateController) TestGetAllTemplate() {
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

			s.mockTemplateService.On("GetAllTemplate", mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.templateController.GetAllTemplate(c)

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

func (s *TestSuiteTemplateController) TestGetTemplateDetail() {
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

			s.mockTemplateService.On("GetTemplateDetail", mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.templateController.GetTemplateDetail(c)

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

func TestTemplateController(t *testing.T) {
	suite.Run(t, new(TestSuiteTemplateController))
}
