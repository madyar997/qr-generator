package main

import (
	"log"

	"github.com/madyar997/maquette/config"
	"github.com/madyar997/maquette/internal/app"
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
