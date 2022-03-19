package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/app"
	"github.com/joho/godotenv"
)

// nolint: nolintlint, gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config/config.yaml", "path to configuration file")
	flag.Parse()

	a, err := app.NewApp(configPath)
	if err != nil {
		log.Fatal(err)
	}

	a.Start()
}
