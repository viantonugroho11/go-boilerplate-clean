package handler

import (
	"strings"

	userEntity "go-boilerplate-clean/internal/entity/users"
	"go-boilerplate-clean/internal/shared/apperrors"
	"go-boilerplate-clean/internal/shared/pagination"
	"go-boilerplate-clean/internal/shared/response"
	"go-boilerplate-clean/internal/transport/apis/dto"
	userUsecase "go-boilerplate-clean/internal/usecase/users"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service userUsecase.UserService
}

func NewUserHandler(service userUsecase.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BindError(c, err)
	}
	user, err := h.service.Create(c.Request().Context(), req.ToEntity())
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, user)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return response.Error(c, apperrors.ErrUserIDRequired)
	}
	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, user)
}

func (h *UserHandler) List(c echo.Context) error {
	page := pagination.ParseQuery(c)
	users, total, err := h.service.List(c.Request().Context(), page)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Paginated(c, pagination.NewList(users, page, total))
}

func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return response.Error(c, apperrors.ErrUserIDRequired)
	}
	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BindError(c, err)
	}
	user, err := h.service.Update(c.Request().Context(), userEntity.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, user)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return response.Error(c, apperrors.ErrUserIDRequired)
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}
	return response.NoContent(c)
}
