package handler

import (
	"net/http"

	"go-boilerplate-clean/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service usecase.UserService
}

func NewUserHandler(service usecase.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) Create(c echo.Context) error {
	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	user, err := h.service.Create(c.Request().Context(), usecase.User{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) List(c echo.Context) error {
	users, err := h.service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req updateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	user, err := h.service.Update(c.Request().Context(), usecase.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
