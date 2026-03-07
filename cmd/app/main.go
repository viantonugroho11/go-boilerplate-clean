package main

import (
	"log"
	"os"

	"go-boilerplate-clean/internal/bootstrap"
)

func main() {
	if err := bootstrap.RunApp(); err != nil {
		log.Fatalf("app: %v", err)
	}
	os.Exit(0)
}
