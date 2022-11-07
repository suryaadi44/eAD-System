package controller

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
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
	api.GET("/documents/:document_id/status", d.GetDocumentStatus)

	secureApi.POST("/templates", d.AddTemplate)
	secureApi.POST("/documents", d.AddDocument)
	secureApi.GET("/documents", d.GetBriefDocument)
	secureApi.GET("/documents/:document_id", d.GetDocument)
	secureApi.GET("/documents/:document_id/pdf", d.GetPDFDocument)
	secureApi.PATCH("/documents/:document_id/verify", d.VerifyDocument)
	secureApi.PATCH("/documents/:document_id/sign", d.SignDocument)
	secureApi.DELETE("/documents/:document_id", d.DeleteDocument)
	secureApi.PUT("/documents/:document_id", d.UpdateDocument)
	secureApi.PUT("/documents/:document_id/fields", d.UpdateDocumentFields)
}

func (d *DocumentController) AddTemplate(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	if role < 2 { // role 2 or above are employee
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	template := new(dto.TemplateRequest)
	if err := c.Bind(template); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(template); err != nil {
		return err
	}

	file, err := c.FormFile("template")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	fileSrc, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer fileSrc.Close()

	err = d.documentService.AddTemplate(c.Request().Context(), template, fileSrc, file.Filename)
	if err != nil {
		if err == utils.ErrDuplicateTemplateName {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
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
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(document); err != nil {
		return err
	}

	claims := d.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	id, err := d.documentService.AddDocument(c.Request().Context(), document, userID)
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
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}
}

func (d *DocumentController) GetBriefDocument(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	userID := claims["user_id"].(string)

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

	documents, err := d.documentService.GetBriefDocuments(c.Request().Context(), userID, int(role), int(pageInt), int(limitInt))
	if err != nil {
		if err == utils.ErrDocumentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting document",
		"data":    documents,
		"meta": echo.Map{
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (d *DocumentController) GetDocumentStatus(c echo.Context) error {
	documentID := c.Param("document_id")
	status, err := d.documentService.GetDocumentStatus(c.Request().Context(), documentID)
	if err != nil {
		if err == utils.ErrDocumentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting document status",
		"data":    status,
	})
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
			return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
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

func (d *DocumentController) VerifyDocument(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	userID := claims["user_id"].(string)
	if role < 2 { // role 2 or above are employee
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	documentID := c.Param("document_id")
	err := d.documentService.VerifyDocument(c.Request().Context(), documentID, userID)
	if err != nil {
		if err == utils.ErrDocumentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success verifying document",
	})
}

func (d *DocumentController) SignDocument(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	userID := claims["user_id"].(string)
	if role < 3 { // role 2 or above are employee
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	documentID := c.Param("document_id")
	err := d.documentService.SignDocument(c.Request().Context(), documentID, userID)
	if err != nil {
		switch err {
		case utils.ErrDocumentNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrNotVerifiedYet:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success signing document",
	})
}

func (d *DocumentController) DeleteDocument(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	userID := claims["user_id"].(string)

	documentID := c.Param("document_id")
	err := d.documentService.DeleteDocument(c.Request().Context(), userID, int(role), documentID)
	if err != nil {
		switch err {
		case utils.ErrDocumentNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrDidntHavePermission:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		case utils.ErrAlreadySigned:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success deleting document",
	})
}

func (d *DocumentController) UpdateDocument(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	if role < 2 {
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	documentID := c.Param("document_id")
	var document dto.DocumentUpdateRequest
	if err := c.Bind(&document); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := d.documentService.UpdateDocument(c.Request().Context(), &document, documentID)
	if err != nil {
		switch err {
		case utils.ErrDocumentNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrAlreadySigned:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		case utils.ErrAlreadyVerified:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success updating document",
	})
}

func (d *DocumentController) UpdateDocumentFields(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	userID := claims["user_id"].(string)

	documentID := c.Param("document_id")
	var fields dto.FieldsUpdateRequest
	if err := c.Bind(&fields); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(fields); err != nil {
		return err
	}

	err := d.documentService.UpdateDocumentFields(c.Request().Context(), userID, int(role), documentID, fields)
	if err != nil {
		switch err {
		case utils.ErrDocumentNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrAlreadySigned:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		case utils.ErrAlreadyVerified:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		case utils.ErrDidntHavePermission:
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success updating document",
	})
}
