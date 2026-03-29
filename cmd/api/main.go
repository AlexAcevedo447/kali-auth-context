package main

import (
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/db"
	"kali-auth-context/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	log, err := logger.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		_ = log.Sync()
	}()

	pool, err := db.NewPool(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer pool.Close()

	app := fiber.New()

	log.Info("Server running on port " + cfg.AppPort)
	err = app.Listen(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err.Error())
	}
}
