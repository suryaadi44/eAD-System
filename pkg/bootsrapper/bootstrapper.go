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
	"github.com/suryaadi44/eAD-System/pkg/utils/jwt_service"
	"github.com/suryaadi44/eAD-System/pkg/utils/password"
	"github.com/suryaadi44/eAD-System/pkg/utils/pdf"
	"github.com/suryaadi44/eAD-System/pkg/utils/qr"
	"github.com/suryaadi44/eAD-System/pkg/utils/validation"

	"time"

	"gorm.io/gorm"
)

func InitController(e *echo.Echo, db *gorm.DB, conf map[string]string) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	jwtService := jwt_service.NewJWTService(conf["JWT_SECRET"], 1*time.Hour)

	v1 := e.Group("/v1")
	secureV1 := v1.Group("")
	secureV1.Use(middleware.JWT([]byte(conf["JWT_SECRET"])))

	qrCodeService := qr.NewCodeServiceImpl(conf["QR_PATH"])
	renderService := renderServicePkg.NewRenderServiceImpl(qrCodeService)
	passwordFunc := password.NewPasswordFuncImpl()
	pdfService := pdf.NewPDFService()

	templateRepository := templateRepositoryPkg.NewTemplateRepositoryImpl(db)
	templateService := templateServicePkg.NewTemplateServiceImpl(templateRepository)
	templateController := templateControllerPkg.NewTemplateController(templateService, jwtService)
	templateController.InitRoute(v1, secureV1)

	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository, passwordFunc, jwtService)
	userController := userControllerPkg.NewUserController(userService, jwtService)
	userController.InitRoute(v1, secureV1)

	documentRepository := documentRepositoryPkg.NewDocumentRepositoryImpl(db)
	documentService := documentServicePkg.NewDocumentServiceImpl(documentRepository, templateRepository, pdfService, renderService)
	documentController := documentControllerPkg.NewDocumentController(documentService, jwtService)
	documentController.InitRoute(v1, secureV1)
}
