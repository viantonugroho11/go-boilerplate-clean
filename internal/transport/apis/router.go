package apis

import (
	"github.com/labstack/echo/v4"
	"go-boilerplate-clean/internal/transport/http/handler"
	"go-boilerplate-clean/internal/usecase"
)

func RegisterRoutes(e *echo.Echo, userService usecase.UserService) {
	userHandler := handler.NewUserHandler(userService)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	users := e.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)
}


