package controller

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"net/http"
	"strconv"
)

type (
	JWTService interface {
		GetClaims(c *echo.Context) jwt.MapClaims
	}

	DocumentController struct {
		documentService service.DocumentService
		jwtService      JWTService
	}
)

func NewDocumentController(documentService service.DocumentService, jwtService JWTService) *DocumentController {
	return &DocumentController{
		documentService: documentService,
		jwtService:      jwtService,
	}
}

func (d *DocumentController) InitRoute(api *echo.Group, secureApi *echo.Group) {
	api.GET("/templates", d.GetAllTemplate)
	api.GET("/templates/:template_id", d.GetTemplateDetail)

	secureApi.POST("/templates", d.AddTemplate)
	secureApi.POST("/documents", d.AddDocument)
	secureApi.GET("/documents/:document_id", d.GetDocument)
	secureApi.GET("/documents/:document_id/pdf", d.GetPDFDocument)
}

func (d *DocumentController) AddTemplate(c echo.Context) error {
	template := new(dto.TemplateRequest)
	if err := c.Bind(template); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody)
	}

	if err := c.Validate(template); err != nil {
		return err
	}

	file, err := c.FormFile("template")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody)
	}

	err = d.documentService.AddTemplate(c.Request().Context(), *template, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success adding template",
	})
}

func (d *DocumentController) GetAllTemplate(c echo.Context) error {
	templates, err := d.documentService.GetAllTemplate(c.Request().Context())
	if err != nil {
		if err == utils.ErrTemplateNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting all template",
		"data":    templates,
	})
}

func (d *DocumentController) GetTemplateDetail(c echo.Context) error {
	templateId := c.Param("template_id")
	templateIdInt, err := strconv.ParseUint(templateId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidTemplateID.Error())
	}

	template, err := d.documentService.GetTemplateDetail(c.Request().Context(), uint(templateIdInt))
	if err != nil {
		if err == utils.ErrTemplateNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting template detail",
		"data":    template,
	})
}

func (d *DocumentController) AddDocument(c echo.Context) error {
	document := new(dto.DocumentRequest)
	if err := c.Bind(document); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody)
	}

	if err := c.Validate(document); err != nil {
		return err
	}

	claims := d.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	id, err := d.documentService.AddDocument(c.Request().Context(), *document, userID)
	if err != nil {
		switch err {
		case utils.ErrTemplateNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrFieldNotMatch:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case utils.ErrDuplicateRegister:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success adding document",
		"data": echo.Map{
			"id": id,
		},
	})
}

func (d *DocumentController) GetDocument(c echo.Context) error {
	documentID := c.Param("document_id")
	document, err := d.documentService.GetDocument(c.Request().Context(), documentID)
	if err != nil {
		if err == utils.ErrDocumentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	claims := d.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)
	role := claims["role"].(float64)

	switch {
	case role > 1:
		fallthrough
	case document.Applicant.ID == userID:
		return c.JSON(http.StatusOK, echo.Map{
			"message": "success getting document",
			"data":    document,
		})
	default:
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDocumentAccessDenied.Error())
	}
}

func (d *DocumentController) GetPDFDocument(c echo.Context) error {
	documentID := c.Param("document_id")

	claims := d.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)
	role := claims["role"].(float64)

	if role == 1 {
		applicantID, err := d.documentService.GetApplicantID(c.Request().Context(), documentID)
		if err != nil {
			if err == utils.ErrDocumentNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if *applicantID != userID {
			return echo.NewHTTPError(http.StatusForbidden, utils.ErrDocumentAccessDenied.Error())
		}
	}

	pdf, err := d.documentService.GeneratePDFDocument(c.Request().Context(), documentID)
	if err != nil {
		if err == utils.ErrDocumentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "application/pdf", pdf)
}
