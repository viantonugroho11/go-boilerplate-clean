package handler

import (
	"net/http"
	"strings"

	"go-boilerplate-clean/internal/transport/apis/dto"
	usecasesample "go-boilerplate-clean/internal/usecase/sample"

	"github.com/labstack/echo/v4"
)

type SampleHandler struct {
	service usecasesample.SampleService
}

func NewSampleHandler(service usecasesample.SampleService) *SampleHandler {
	return &SampleHandler{service: service}
}

func (h *SampleHandler) Create(c echo.Context) error {
	var req dto.SaveSampleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	sample, err := h.service.Save(c.Request().Context(), req.ToEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, sample)
}

func (h *SampleHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	var req dto.SaveSampleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	entity := req.ToEntity()
	entity.ID = id
	sample, err := h.service.Save(c.Request().Context(), entity)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, sample)
}
