package apis

import (
	"github.com/labstack/echo/v5"

	"go-boilerplate-clean/internal/transport/apis/handler"
	usecasesample "go-boilerplate-clean/internal/usecase/sample"
	"go-boilerplate-clean/internal/usecase/users"
)

func RegisterRoutes(e *echo.Echo, userService users.UserService, sampleService usecasesample.SampleService) {
	userHandler := handler.NewUserHandler(userService)
	sampleHandler := handler.NewSampleHandler(sampleService)

	e.GET("/healthz", func(c *echo.Context) error {
		return c.String(200, "ok")
	})

	users := e.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	samples := e.Group("/samples")
	samples.POST("", sampleHandler.Create)
	samples.PUT("/:id", sampleHandler.Update)
}
