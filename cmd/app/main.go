package main

import (
	"log"

	"github.com/madyar997/qr-generator/config"
	"github.com/madyar997/qr-generator/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
