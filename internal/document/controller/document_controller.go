package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/document/dto"
	"github.com/suryaadi44/eAD-System/internal/document/service"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"net/http"
)

type DocumentController struct {
	documentService service.DocumentService
}

func NewDocumentController(documentService service.DocumentService) *DocumentController {
	return &DocumentController{documentService}
}

func (d *DocumentController) InitRoute(api *echo.Group, secureApi *echo.Group) {
	api.GET("/templates", d.GetAllTemplate)
	secureApi.POST("/templates", d.AddTemplate)
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
