package main

import (
	stdlog "log"

	bootstrap "kali-auth-context/internal/bootstrap/di"
	"kali-auth-context/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	container, err := bootstrap.InitializeContainer()
	if err != nil {
		stdlog.Fatalf("failed to initialize container: %v", err)
	}
	defer func() {
		if closeErr := container.Close(); closeErr != nil {
			stdlog.Printf("failed to close container resources: %v", closeErr)
		}
	}()

	appLogger, err := logger.New()
	if err != nil {
		stdlog.Fatalf("failed to initialize logger: %v", err)
	}

	defer func() {
		_ = appLogger.Sync()
	}()

	app := fiber.New()
	container.HTTP.Router.Register(app)

	appLogger.Info("Server running on port " + container.Config.AppPort)
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		appLogger.Fatal(err.Error())
	}
}
