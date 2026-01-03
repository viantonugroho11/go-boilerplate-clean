package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-boilerplate-clean/internal/config"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	redisinfra "go-boilerplate-clean/internal/infrastructure/cache/redis"
	pginfra "go-boilerplate-clean/internal/infrastructure/database/postgres"
	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/apis"
	kafkarunner "go-boilerplate-clean/internal/transport/event/kafka"
	"go-boilerplate-clean/internal/usecase"

	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Wiring dependencies
	ctx := context.Background()
	db, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err := pginfra.Migrate(db); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	userRepo := userpg.NewUserRepository(db)
	userService := usecase.NewUserService(userRepo)
	apis.RegisterRoutes(e, userService)

	// Init Redis
	redisClient, err := redisinfra.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("redis init error: %v", err)
	}
	defer redisClient.Close()

	// Init Kafka Producer
	producer, err := kafkainfra.NewProducer(cfg.KafkaBrokersList(), cfg.KafkaClientID)
	if err != nil {
		log.Fatalf("kafka producer init error: %v", err)
	}
	defer producer.Close()

	// Init Kafka Consumer
	consumerHandler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		return kafkarunner.ExampleHandler(ctx, msg.Key, msg.Value)
	}
	consumer, err := kafkainfra.NewConsumer(cfg.KafkaBrokersList(), cfg.KafkaGroupID, cfg.KafkaTopic, consumerHandler)
	if err != nil {
		log.Fatalf("kafka consumer init error: %v", err)
	}
	kafkarunner.RegisterConsumers(ctx, consumer)
	defer consumer.Close()

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
