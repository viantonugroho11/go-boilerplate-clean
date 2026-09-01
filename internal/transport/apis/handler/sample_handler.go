package handler

import (
	"strings"

	"go-boilerplate-clean/internal/shared/apperrors"
	"go-boilerplate-clean/internal/shared/response"
	"go-boilerplate-clean/internal/transport/apis/dto"
	usecasesample "go-boilerplate-clean/internal/usecase/sample"

	"github.com/labstack/echo/v5"
)

type SampleHandler struct {
	service usecasesample.SampleService
}

func NewSampleHandler(service usecasesample.SampleService) *SampleHandler {
	return &SampleHandler{service: service}
}

func (h *SampleHandler) Create(c *echo.Context) error {
	var req dto.SaveSampleRequest
	if err := c.Bind(&req); err != nil {
		return response.BindError(c, err)
	}
	sample, err := h.service.Save(c.Request().Context(), req.ToEntity())
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, sample)
}

func (h *SampleHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return response.Error(c, apperrors.ErrSampleIDRequired)
	}
	var req dto.SaveSampleRequest
	if err := c.Bind(&req); err != nil {
		return response.BindError(c, err)
	}
	entity := req.ToEntity()
	entity.ID = id
	sample, err := h.service.Save(c.Request().Context(), entity)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, sample)
}
