package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/user/dto"
	"github.com/suryaadi44/eAD-System/internal/user/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) InitRoute(e *echo.Group) {
	e.POST("/signup", u.SignUpUser)
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

	return c.JSON(200, echo.Map{
		"message": "success creating user",
	})
}
