package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
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

func (m *MockUserService) GetBriefUsers(ctx context.Context, page int, limit int) (*dto.BriefUsersResponse, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(*dto.BriefUsersResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID string, request *dto.UserUpdateRequest) error {
	args := m.Called(ctx, userID, request)
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

type TestSuiteUserControllers struct {
	suite.Suite
	mockUserService *MockUserService
	mockValidator   *MockValidator
	mockJWT         *MockJWTService
	userController  *UserController
	echoApp         *echo.Echo
}

func (s *TestSuiteUserControllers) SetupTest() {
	s.mockUserService = new(MockUserService)
	s.mockValidator = new(MockValidator)
	s.mockJWT = new(MockJWTService)
	s.userController = NewUserController(s.mockUserService, s.mockJWT)
	s.echoApp = echo.New()
	s.echoApp.Validator = s.mockValidator
}

func (s *TestSuiteUserControllers) TearDownTest() {
	s.mockUserService = nil
	s.userController = nil
	s.echoApp = nil
}

func (s *TestSuiteUserControllers) TestInitRoute() {
	group := s.echoApp.Group("/user")
	s.NotPanics(func() {
		s.userController.InitRoute(group, group)
	})
}

func (s *TestSuiteUserControllers) TestUpdateUser() {
	for _, tc := range []struct {
		Name            string
		RequestBody     interface{}
		FunctionError   error
		ValidationError error
		JWTReturn       jwt.MapClaims
		ExpectedStatus  int
		ExpectedBody    echo.Map
		ExpectedError   error
	}{
		{
			Name: "Success updating user",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
			},
			FunctionError:   nil,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success update user",
			},
			ExpectedError: nil,
		},
		{
			Name: "Failed updating user: user not found",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
			},
			FunctionError:   utils.ErrUserNotFound,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrUserNotFound,
		},
		{
			Name: "Failed updating user : Username already exist",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
			},
			FunctionError:   utils.ErrUsernameAlreadyExist,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrUsernameAlreadyExist,
		},
		{
			Name: "Failed updating user : NIK already exist",
			RequestBody: &dto.UserSignUpRequest{
				NIK: "1234567890123456",
			},
			FunctionError:   utils.ErrNIKAlreadyExist,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrNIKAlreadyExist,
		},
		{
			Name: "Failed updating user : NIP already exist",
			RequestBody: &dto.UserSignUpRequest{
				NIP: "123456789012345678",
			},
			FunctionError:   utils.ErrNIPAlreadyExist,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrNIPAlreadyExist,
		},
		{
			Name: "Failed updating user : Generic error",
			RequestBody: &dto.UserSignUpRequest{
				Username: "suryaadi",
			},
			FunctionError:   errors.New("generic error"),
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("generic error"),
		},
		{
			Name:            "Failed updating user : Invalid request body",
			RequestBody:     "invalid request body",
			FunctionError:   nil,
			ValidationError: nil,
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrBadRequestBody,
		},
		{
			Name:            "Failed updating user : Validation error",
			RequestBody:     &dto.UserSignUpRequest{},
			FunctionError:   nil,
			ValidationError: echo.NewHTTPError(http.StatusBadRequest, "validation error"),
			JWTReturn: jwt.MapClaims{
				"user_id": "1",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("validation error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoApp.NewContext(r, w)

			s.mockJWT.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockValidator.On("Validate", mock.Anything).Return(tc.ValidationError)
			s.mockUserService.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionError)

			err = s.userController.UpdateUser(c)

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
				s.mockValidator.On("Validate", tc.RequestBody).Return(tc.ValidationError)
			} else {
				s.mockValidator.On("Validate", tc.RequestBody).Return(nil)
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

func (s *TestSuiteUserControllers) TestGetBriefUsers() {
	for _, tc := range []struct {
		Name           string
		Page           string
		Limit          string
		FunctionError  error
		FunctionReturn *dto.BriefUsersResponse
		JWTReturn      jwt.MapClaims
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name:          "Success getting brief users",
			Page:          "1",
			Limit:         "10",
			FunctionError: nil,
			FunctionReturn: &dto.BriefUsersResponse{
				{
					ID:       "1",
					Username: "user",
					Name:     "user",
				},
			},
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success get users",
				"data": []interface{}{
					map[string]interface{}{
						"id":       "1",
						"username": "user",
						"name":     "user",
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
			Name:          "Success getting brief users: blank page and limit",
			Page:          "",
			Limit:         "",
			FunctionError: nil,
			FunctionReturn: &dto.BriefUsersResponse{
				{
					ID:       "1",
					Username: "user",
					Name:     "user",
				},
			},
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: echo.Map{
				"message": "success get users",
				"data": []interface{}{
					map[string]interface{}{
						"id":       "1",
						"username": "user",
						"name":     "user",
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
			Name:           "Failed getting brief users : Invalid page",
			Page:           "invalid",
			Limit:          "10",
			FunctionError:  nil,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrInvalidNumber,
		},
		{
			Name:           "Failed getting brief users : Invalid limit",
			Page:           "1",
			Limit:          "invalid",
			FunctionError:  nil,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrInvalidNumber,
		},
		{
			Name:           "Failed getting brief users : role is not employee",
			Page:           "1",
			Limit:          "10",
			FunctionError:  nil,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(1),
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrDidntHavePermission,
		},
		{
			Name:           "Failed getting brief users : error no user found",
			Page:           "",
			Limit:          "",
			FunctionError:  utils.ErrUserNotFound,
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   nil,
			ExpectedError:  utils.ErrUserNotFound,
		},
		{
			Name:           "Failed getting brief users : error from service",
			Page:           "",
			Limit:          "",
			FunctionError:  errors.New("error from service"),
			FunctionReturn: nil,
			JWTReturn: jwt.MapClaims{
				"role": float64(2),
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   nil,
			ExpectedError:  errors.New("error from service"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			r := httptest.NewRequest("GET", "/users", nil)

			w := httptest.NewRecorder()

			q := r.URL.Query()
			q.Add("page", tc.Page)
			q.Add("limit", tc.Limit)
			r.URL.RawQuery = q.Encode()

			c := s.echoApp.NewContext(r, w)

			s.mockJWT.On("GetClaims", mock.Anything).Return(tc.JWTReturn)
			s.mockUserService.On("GetBriefUsers", mock.Anything, mock.Anything, mock.Anything).Return(tc.FunctionReturn, tc.FunctionError)

			err := s.userController.GetBriefUsers(c)

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
				s.mockValidator.On("Validate", tc.RequestBody).Return(tc.ValidationError)
			} else {
				s.mockValidator.On("Validate", tc.RequestBody).Return(nil)
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

func TestUserControllers(t *testing.T) {
	suite.Run(t, new(TestSuiteUserControllers))
}
