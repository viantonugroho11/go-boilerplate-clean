package bootstrap

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/event"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"
)

const (
	ConsumerUser  = event.ConsumerNameUser
	ConsumerOrder = event.ConsumerNameOrder
)

// RunConsumer menjalankan consumer sesuai name (user | order): config global, wiring terisolasi per consumer, run sampai signal.
func RunConsumer(name string) error {
	if err := LoadConfig(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := Config()
	var consumer interface{ Close() error }
	switch name {
	case ConsumerUser:
		db, err := initDB()
		if err != nil {
			return err
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		userRepo := userpg.NewUserRepository(db)
		userService := usecaseusers.NewUserService(userRepo)
		c, err := event.RunUser(ctx, cfg, userService)
		if err != nil {
			return err
		}
		consumer = c
	case ConsumerOrder:
		c, err := event.RunOrder(ctx, cfg)
		if err != nil {
			return err
		}
		consumer = c
	default:
		return fmt.Errorf("consumer %q tidak dikenal, pilih: %s", name, strings.Join(event.ConsumerNames(), " | "))
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("consumer close error: %v", err)
		}
	}()

	log.Printf("running consumer: %s", name)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, stopping consumer...")
	return nil
}

// ParseConsumerFlag parse flag -consumer= dan validasi. Return nama consumer atau empty.
func ParseConsumerFlag() string {
	consumerFlag := flag.String("consumer", "", "nama consumer (user | order)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -consumer=<name>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  -consumer string   pilih: %s | %s\n", ConsumerUser, ConsumerOrder)
		flag.PrintDefaults()
	}
	flag.Parse()
	return *consumerFlag
}
