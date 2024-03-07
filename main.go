package main

import (
	"fmt"
	"log/slog"
	"service/app/routes"
	"service/pkg/configs"
	"service/pkg/logging"

	"github.com/gofiber/fiber/v2"
)



func main() {



	cfg := configs.GetConfig()
	


	log := logging.SetupLogger(cfg.ModeLog)
	log.Info("Starting service", slog.String("env", cfg.ModeLog))

	app := fiber.New()

	routes.Init(app, cfg, log)

	address := fmt.Sprintf("%s:%s", cfg.Fiber.BindIp, cfg.Fiber.Port)

	if err := app.Listen(address); err != nil {
		fmt.Printf("server not running %v", err)
	}
}

