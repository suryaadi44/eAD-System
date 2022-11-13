package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	documentControllerPkg "github.com/suryaadi44/eAD-System/internal/document/controller"
	templateControllerPkg "github.com/suryaadi44/eAD-System/internal/template/controller"
	userControllerPkg "github.com/suryaadi44/eAD-System/internal/user/controller"
	"github.com/suryaadi44/eAD-System/pkg/utils/validation"
)

type Routes struct {
	userController     *userControllerPkg.UserController
	templateController *templateControllerPkg.TemplateController
	documentController *documentControllerPkg.DocumentController
}

func NewRoutes(userController *userControllerPkg.UserController, templateController *templateControllerPkg.TemplateController, documentController *documentControllerPkg.DocumentController) *Routes {
	return &Routes{
		userController:     userController,
		templateController: templateController,
		documentController: documentController,
	}
}

func (r *Routes) Init(e *echo.Echo, conf map[string]string) {
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Recover())

	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(conf["JWT_SECRET"]),
	})

	v1 := e.Group("/v1")

	// Users
	users := v1.Group("/users")
	users.POST("/signup/", r.userController.SignUpUser)
	users.POST("/login/", r.userController.LoginUser)

	usersWithAuth := users.Group("", jwtMiddleware)
	usersWithAuth.GET("/", r.userController.GetBriefUsers)
	usersWithAuth.PUT("/", r.userController.UpdateUser)

	// Documents
	documents := v1.Group("/documents")
	documents.GET("/:document_id/status/", r.documentController.GetDocumentStatus)

	documentsWithAuth := documents.Group("", jwtMiddleware)
	documentsWithAuth.POST("/", r.documentController.AddDocument)
	documentsWithAuth.GET("/", r.documentController.GetBriefDocument)
	documentsWithAuth.GET("/:document_id/", r.documentController.GetDocument)
	documentsWithAuth.GET("/:document_id/pdf/", r.documentController.GetPDFDocument)
	documentsWithAuth.PATCH("/:document_id/verify/", r.documentController.VerifyDocument)
	documentsWithAuth.PATCH("/:document_id/sign/", r.documentController.SignDocument)
	documentsWithAuth.DELETE("/:document_id/", r.documentController.DeleteDocument)
	documentsWithAuth.PUT("/:document_id/", r.documentController.UpdateDocument)
	documentsWithAuth.PUT("/:document_id/fields/", r.documentController.UpdateDocumentFields)

	// Templates
	templates := v1.Group("/templates")
	templates.GET("/", r.templateController.GetAllTemplate)
	templates.GET("/:template_id/", r.templateController.GetTemplateDetail)

	templatesWithAuth := templates.Group("", jwtMiddleware)
	templatesWithAuth.POST("/", r.templateController.AddTemplate)
}
