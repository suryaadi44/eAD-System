package controller

import (
	error2 "github.com/suryaadi44/eAD-System/pkg/utils/error"
	"github.com/suryaadi44/eAD-System/pkg/utils/jwt_service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suryaadi44/eAD-System/internal/template/dto"
	"github.com/suryaadi44/eAD-System/internal/template/service"
)

type TemplateController struct {
	templateService service.TemplateService
	jwtService      jwt_service.JWTService
}

func NewTemplateController(templateService service.TemplateService, jwtService jwt_service.JWTService) *TemplateController {
	return &TemplateController{
		templateService: templateService,
		jwtService:      jwtService,
	}
}

func (t *TemplateController) InitRoute(api *echo.Group, secureApi *echo.Group) {
	api.GET("/templates", t.GetAllTemplate)
	api.GET("/templates/:template_id", t.GetTemplateDetail)

	secureApi.POST("/templates", t.AddTemplate)
}

func (t *TemplateController) AddTemplate(c echo.Context) error {
	claims := t.jwtService.GetClaims(&c)
	role := claims["role"].(float64)
	if role < 2 { // role 2 or above are employee
		return echo.NewHTTPError(http.StatusForbidden, error2.ErrDidntHavePermission.Error())
	}

	template := new(dto.TemplateRequest)
	if err := c.Bind(template); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, error2.ErrBadRequestBody.Error())
	}

	if err := c.Validate(template); err != nil {
		return err
	}

	file, err := c.FormFile("template")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, error2.ErrBadRequestBody.Error())
	}

	fileSrc, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer fileSrc.Close()

	err = t.templateService.AddTemplate(c.Request().Context(), template, fileSrc, file.Filename)
	if err != nil {
		if err == error2.ErrDuplicateTemplateName {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success adding template",
	})
}

func (t *TemplateController) GetAllTemplate(c echo.Context) error {
	templates, err := t.templateService.GetAllTemplate(c.Request().Context())
	if err != nil {
		if err == error2.ErrTemplateNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting all template",
		"data":    templates,
	})
}

func (t *TemplateController) GetTemplateDetail(c echo.Context) error {
	templateId := c.Param("template_id")
	templateIdInt, err := strconv.ParseUint(templateId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, error2.ErrInvalidTemplateID.Error())
	}

	template, err := t.templateService.GetTemplateDetail(c.Request().Context(), uint(templateIdInt))
	if err != nil {
		if err == error2.ErrTemplateNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting template detail",
		"data":    template,
	})
}
