package main

import (
	"flag"
	"log"
	"os"

	"go-boilerplate-clean/internal/bootstrap"
)

func main() {
	name := bootstrap.ParseConsumerFlag()
	if name == "" {
		log.Printf("flag -consumer wajib diisi. contoh: -consumer=user atau -consumer=order")
		flag.Usage()
		os.Exit(1)
	}

	if err := bootstrap.RunConsumer(name); err != nil {
		log.Fatalf("consumer: %v", err)
	}
	os.Exit(0)
}
