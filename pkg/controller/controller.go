package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	userControllerPkg "github.com/suryaadi44/eAD-System/internal/user/controller"
	userRepositoryPkg "github.com/suryaadi44/eAD-System/internal/user/repository"
	userServicePkg "github.com/suryaadi44/eAD-System/internal/user/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
	"time"
)

func InitController(e *echo.Echo, db *gorm.DB, conf map[string]string) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	jwtService := utils.NewJWTService(conf["JWT_SECRET"], 1*time.Hour)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository, utils.PasswordFunc{}, jwtService)
	userController := userControllerPkg.NewUserController(userService)
	userController.InitRoute(v1)
}
