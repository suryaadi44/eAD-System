package bootsrapper

import (
	"github.com/labstack/echo/v4"
	documentControllerPkg "github.com/suryaadi44/eAD-System/internal/document/controller"
	documentRepositoryPkg "github.com/suryaadi44/eAD-System/internal/document/repository/impl"
	documentServicePkg "github.com/suryaadi44/eAD-System/internal/document/service/impl"
	templateControllerPkg "github.com/suryaadi44/eAD-System/internal/template/controller"
	templateRepositoryPkg "github.com/suryaadi44/eAD-System/internal/template/repository/impl"
	templateServicePkg "github.com/suryaadi44/eAD-System/internal/template/service/impl"
	userControllerPkg "github.com/suryaadi44/eAD-System/internal/user/controller"
	userRepositoryPkg "github.com/suryaadi44/eAD-System/internal/user/repository/impl"
	userServicePkg "github.com/suryaadi44/eAD-System/internal/user/service/impl"
	"github.com/suryaadi44/eAD-System/pkg/routes"
	renderServicePkg "github.com/suryaadi44/eAD-System/pkg/utils/html/impl"
	jwtPkg "github.com/suryaadi44/eAD-System/pkg/utils/jwt_service/impl"
	passwordPkg "github.com/suryaadi44/eAD-System/pkg/utils/password/impl"
	pdfPkg "github.com/suryaadi44/eAD-System/pkg/utils/pdf/impl"
	qrPkg "github.com/suryaadi44/eAD-System/pkg/utils/qr/impl"
	"time"

	"gorm.io/gorm"
)

func InitController(e *echo.Echo, db *gorm.DB, conf map[string]string) {
	qrCodeService := qrPkg.NewCodeServiceImpl(conf["QR_PATH"])
	renderService := renderServicePkg.NewRenderServiceImpl(qrCodeService)
	passwordFunc := passwordPkg.NewPasswordFuncImpl()
	pdfService := pdfPkg.NewPDFService()
	jwtService := jwtPkg.NewJWTService(conf["JWT_SECRET"], 1*time.Hour)

	// User
	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository, passwordFunc, jwtService)
	userController := userControllerPkg.NewUserController(userService, jwtService)

	// Template
	templateRepository := templateRepositoryPkg.NewTemplateRepositoryImpl(db)
	templateService := templateServicePkg.NewTemplateServiceImpl(templateRepository)
	templateController := templateControllerPkg.NewTemplateController(templateService, jwtService)

	// Document
	documentRepository := documentRepositoryPkg.NewDocumentRepositoryImpl(db)
	documentService := documentServicePkg.NewDocumentServiceImpl(documentRepository, templateRepository, pdfService, renderService)
	documentController := documentControllerPkg.NewDocumentController(documentService, jwtService)

	route := routes.NewRoutes(userController, templateController, documentController)
	route.Init(e, conf)
}
