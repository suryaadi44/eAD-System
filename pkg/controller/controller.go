package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func InitController(e *echo.Echo, _ *gorm.DB) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

}
