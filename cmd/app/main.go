package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-boilerplate-clean/internal/config"
	pginfra "go-boilerplate-clean/internal/infrastructure/database/postgres"
	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/apis"
	"go-boilerplate-clean/internal/usecase"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Wiring dependencies
	ctx := context.Background()
	pool, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()
	if err := pginfra.Migrate(ctx, pool); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	userRepo := userpg.NewUserRepository(pool)
	userService := usecase.NewUserService(userRepo)
	apis.RegisterRoutes(e, userService)

	// HTTP server with graceful shutdown
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	log.Printf("server listening on :%s", cfg.Port)

	// wait for interrupt signal to gracefully shutdown the server with a timeout
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown gracefully")
	}
}


