package main

import (
	"fmt"
	"log/slog"
	"os"
	"service/app/routes"
	"service/pkg/configs"
	"service/pkg/logging"

	"github.com/gofiber/fiber/v2"
)



func main() {

	// os.Setenv("CONFIG_PATH", "./config.yml")
	// os.Setenv("SECRET_KEY", "Medods_Task1")
	cfg := configs.GetConfig()
	
	fmt.Print(cfg)

	log := logging.SetupLogger(cfg.ModeLog)
	log.Info("Starting service", slog.String("env", cfg.ModeLog))

	app := fiber.New()

	routes.Init(app, cfg, log)

	address := fmt.Sprintf("%s:%s", cfg.Fiber.BindIp, cfg.Fiber.Port)

	if err := app.Listen(address); err != nil {
		fmt.Printf("server not running %v", err)
	}
}

