package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		validationErr := err.(validator.ValidationErrors)
		for _, each := range validationErr {
			switch each.Tag() {
			case "required":
				msg := fmt.Sprintf("%s is required", each.Field())
				return echo.NewHTTPError(http.StatusBadRequest, msg)
			case "len":
				msg := fmt.Sprintf("%s must be %s characters long", each.Field(), each.Param())
				return echo.NewHTTPError(http.StatusBadRequest, msg)
			default:
				msg := fmt.Sprintf("Invalid field %s", each.Field())
				return echo.NewHTTPError(http.StatusBadRequest, msg)
			}
		}
	}

	return nil
}
