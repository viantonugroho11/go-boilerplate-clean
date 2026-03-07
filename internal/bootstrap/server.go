package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-boilerplate-clean/internal/transport/apis"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewEcho buat Echo, middleware, dan daftar routes.
func newEcho(userService usecaseusers.UserService) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover(), middleware.Logger())
	apis.RegisterRoutes(e, userService)
	return e
}

// RunHTTP jalankan server sampai dapat signal interrupt, lalu graceful shutdown. Pakai Config() global untuk port.
func runHTTP(e *echo.Echo) error {
	c := Config()
	server := &http.Server{
		Addr:         ":" + c.App.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()
	log.Printf("server listening on :%s", c.App.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
		return err
	}
	log.Println("server shutdown gracefully")
	return nil
}
