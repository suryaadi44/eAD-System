package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUpUser(ctx context.Context, user *dto.UserSignUpRequest) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserService) LogInUser(ctx context.Context, user *dto.UserLoginRequest) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(a0 interface{}) error {
	args := m.Called(a0)
	return args.Error(0)
}

type TestSuiteUserControllers struct {
	suite.Suite
	mockUserService *MockUserService
	MockValidator   *MockValidator
	userController  *UserController
	echoApp         *echo.Echo
}

func (s *TestSuiteUserControllers) SetupTest() {
	s.mockUserService = new(MockUserService)
	s.MockValidator = new(MockValidator)
	s.userController = NewUserController(s.mockUserService)
	s.echoApp = echo.New()
	s.echoApp.Validator = s.MockValidator
}

func (s *TestSuiteUserControllers) TearDownTest() {
	s.mockUserService = nil
	s.userController = nil
	s.echoApp = nil
}

func (s *TestSuiteUserControllers) TestInitRoute() {
	group := s.echoApp.Group("/user")
	s.NotPanics(func() {
		s.userController.InitRoute(group)
	})
}

func (s *TestSuiteUserControllers) TestSignUpUser() {
	for _, tc := range []struct {
		Name            string
		RequestBody     interface{}
		FunctionError   error
		ValidationError error
		ExpectedStatus  int
		ExpectedBody    echo.Map
		ExpectedError   error
	}{
		{
			Name: "Success creating user",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
				Password: "123456",
				NIK:      "1234567890123456",
				NIP:      "123456789012345678",
				Name:     "Surya Adi",
				Telp:     "081234567890",
				Sex:      "L",
				Address:  "Jl. Jalan",
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedBody: echo.Map{
				"message": "success creating user",
			},
		},
		{
			Name: "Failed creating user : Username already exist",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
				Password: "123456",
				NIK:      "1234567890123456",
				NIP:      "123456789012345678",
				Name:     "Surya Adi",
				Telp:     "081234567890",
				Sex:      "L",
				Address:  "Jl. Jalan",
			},
			FunctionError:  utils.ErrUsernameAlreadyExist,
			ExpectedStatus: http.StatusConflict,
			ExpectedError:  utils.ErrUsernameAlreadyExist,
		},
		{
			Name: "Failed creating user : NIK already exist",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
				Password: "123456",
				NIK:      "1234567890123456",
				NIP:      "123456789012345678",
				Name:     "Surya Adi",
				Telp:     "081234567890",
				Sex:      "L",
				Address:  "Jl. Jalan",
			},
			FunctionError:  utils.ErrNIKAlreadyExist,
			ExpectedStatus: http.StatusConflict,
			ExpectedError:  utils.ErrNIKAlreadyExist,
		},
		{
			Name: "Failed creating user : NIP already exist",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
				Password: "123456",
				NIK:      "1234567890123456",
				NIP:      "123456789012345678",
				Name:     "Surya Adi",
				Telp:     "081234567890",
				Sex:      "L",
				Address:  "Jl. Jalan",
			},
			FunctionError:  utils.ErrNIPAlreadyExist,
			ExpectedStatus: http.StatusConflict,
			ExpectedError:  utils.ErrNIPAlreadyExist,
		},
		{
			Name: "Failed creating user : Generic error",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
				Password: "123456",
				NIK:      "1234567890123456",
				NIP:      "123456789012345678",
				Name:     "Surya Adi",
				Telp:     "081234567890",
				Sex:      "L",
				Address:  "Jl. Jalan",
			},
			FunctionError:  errors.New("generic error"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:           "Failed creating user : Invalid request body",
			RequestBody:    "invalid request body",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name:            "Failed creating user : Validation error",
			RequestBody:     &dto.UserSignUpRequest{},
			ValidationError: echo.NewHTTPError(http.StatusBadRequest, "validation error"),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedError:   errors.New("validation error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)

			s.mockUserService.On("SignUpUser", mock.Anything, tc.RequestBody).Return(tc.FunctionError)

			if tc.ValidationError != nil {
				s.MockValidator.On("Validate", tc.RequestBody).Return(tc.ValidationError)
			} else {
				s.MockValidator.On("Validate", tc.RequestBody).Return(nil)
			}

			err = s.userController.SignUpUser(c)

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

func (s *TestSuiteUserControllers) TestLogInUser() {
	for _, tc := range []struct {
		Name            string
		RequestBody     interface{}
		FunctionError   error
		FunctionReturn  string
		ValidationError error
		ExpectedStatus  int
		ExpectedBody    echo.Map
		ExpectedError   error
	}{
		{
			Name: "Success creating user",
			RequestBody: &dto.UserLoginRequest{
				Username: "suryaadi",
				Password: "123456",
			},
			FunctionReturn: "token",
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success login",
				"token":   "token",
			},
		},
		{
			Name: "Failed creating user : Invalid Credentials",
			RequestBody: &dto.UserLoginRequest{
				Username: "suryaadi",
				Password: "123456",
			},
			FunctionError:  utils.ErrInvalidCredentials,
			ExpectedStatus: http.StatusUnauthorized,
			ExpectedError:  utils.ErrInvalidCredentials,
		},
		{
			Name: "Failed creating user : Generic error",
			RequestBody: &dto.UserLoginRequest{
				Username: "suryaadi",
				Password: "123456",
			},
			FunctionError:  errors.New("generic error"),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:           "Failed creating user : Invalid request body",
			RequestBody:    "invalid request body",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name:            "Failed logging in user : Validation error",
			RequestBody:     &dto.UserLoginRequest{},
			ValidationError: echo.NewHTTPError(http.StatusBadRequest, "validation error"),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedError:   errors.New("validation error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)

			s.mockUserService.On("LogInUser", mock.Anything, tc.RequestBody).Return(tc.FunctionReturn, tc.FunctionError)

			if tc.ValidationError != nil {
				s.MockValidator.On("Validate", tc.RequestBody).Return(tc.ValidationError)
			} else {
				s.MockValidator.On("Validate", tc.RequestBody).Return(nil)
			}

			err = s.userController.LoginUser(c)

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

func TestUserControllers(t *testing.T) {
	suite.Run(t, new(TestSuiteUserControllers))
}
