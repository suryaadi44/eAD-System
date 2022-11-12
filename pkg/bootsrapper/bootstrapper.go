package bootsrapper

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	documentControllerPkg "github.com/suryaadi44/eAD-System/internal/document/controller"
	documentRepositoryPkg "github.com/suryaadi44/eAD-System/internal/document/repository"
	documentServicePkg "github.com/suryaadi44/eAD-System/internal/document/service"
	templateControllerPkg "github.com/suryaadi44/eAD-System/internal/template/controller"
	templateRepositoryPkg "github.com/suryaadi44/eAD-System/internal/template/repository"
	templateServicePkg "github.com/suryaadi44/eAD-System/internal/template/service"
	userControllerPkg "github.com/suryaadi44/eAD-System/internal/user/controller"
	userRepositoryPkg "github.com/suryaadi44/eAD-System/internal/user/repository"
	userServicePkg "github.com/suryaadi44/eAD-System/internal/user/service"
	renderServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/html"
	pdfServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/pdf"
	qrServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/qr"

	"time"

	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/gorm"
)

func InitController(e *echo.Echo, db *gorm.DB, conf map[string]string) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	jwtService := utils.NewJWTService(conf["JWT_SECRET"], 1*time.Hour)

	v1 := e.Group("/v1")
	secureV1 := v1.Group("")
	secureV1.Use(middleware.JWT([]byte(conf["JWT_SECRET"])))

	qrCodeService := qrServicePkg.NewCodeServiceImpl(conf["QR_PATH"])
	renderService := renderServicePkg.NewRenderServiceImpl(qrCodeService)

	templateRepository := templateRepositoryPkg.NewTemplateRepositoryImpl(db)
	templateService := templateServicePkg.NewTemplateServiceImpl(templateRepository)
	templateController := templateControllerPkg.NewTemplateController(templateService, jwtService)
	templateController.InitRoute(v1, secureV1)

	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository, utils.PasswordFunc{}, jwtService)
	userController := userControllerPkg.NewUserController(userService, jwtService)
	userController.InitRoute(v1, secureV1)

	pdfService := pdfServicePkg.NewPDFService()
	documentRepository := documentRepositoryPkg.NewDocumentRepositoryImpl(db)
	documentService := documentServicePkg.NewDocumentServiceImpl(documentRepository, templateRepository, pdfService, renderService)
	documentController := documentControllerPkg.NewDocumentController(documentService, jwtService)
	documentController.InitRoute(v1, secureV1)
}
