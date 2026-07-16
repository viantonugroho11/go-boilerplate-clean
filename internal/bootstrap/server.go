package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-boilerplate-clean/internal/config"
	"go-boilerplate-clean/internal/transport/apis"
	usecasesample "go-boilerplate-clean/internal/usecase/sample"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewEcho buat Echo, middleware, dan daftar routes.
func newEcho(userService usecaseusers.UserService, sampleService usecasesample.SampleService) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover(), middleware.Logger())
	apis.RegisterRoutes(e, userService, sampleService)
	return e
}

func runHTTP(cfg *config.Configuration, e *echo.Echo) error {
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()
	slog.Info("server listening", "port", cfg.App.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "err", err)
		return err
	}
	slog.Info("server shutdown gracefully")
	return nil
}
