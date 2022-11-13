package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/internal/user/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"github.com/suryaadi44/eAD-System/pkg/utils/jwt_service"
	"net/http"
	"strconv"
)

type UserController struct {
	userService service.UserService
	jwtService  jwt_service.JWTService
}

func NewUserController(userService service.UserService, jwtService jwt_service.JWTService) *UserController {
	return &UserController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (u *UserController) SignUpUser(c echo.Context) error {
	user := new(dto.UserSignUpRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := u.userService.SignUpUser(c.Request().Context(), user)
	if err != nil {
		switch err {
		case utils.ErrUsernameAlreadyExist:
			fallthrough
		case utils.ErrNIKAlreadyExist:
			fallthrough
		case utils.ErrNIPAlreadyExist:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "success creating user",
	})
}

func (u *UserController) LoginUser(c echo.Context) error {
	user := new(dto.UserLoginRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	token, err := u.userService.LogInUser(c.Request().Context(), user)
	if err != nil {
		switch err {
		case utils.ErrInvalidCredentials:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success login",
		"token":   token,
	})
}

func (u *UserController) GetBriefUsers(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role"].(float64)

	if role == 1 {
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	users, err := u.userService.GetBriefUsers(c.Request().Context(), int(pageInt), int(limitInt))
	if err != nil {
		if err == utils.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get users",
		"data":    users,
		"meta": echo.Map{
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (u *UserController) UpdateUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	user := new(dto.UserUpdateRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := u.userService.UpdateUser(c.Request().Context(), userID, user)
	if err != nil {
		switch err {
		case utils.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrUsernameAlreadyExist:
			fallthrough
		case utils.ErrNIKAlreadyExist:
			fallthrough
		case utils.ErrNIPAlreadyExist:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update user",
	})
}
