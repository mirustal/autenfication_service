package routes

import (
	"context"
	"log/slog"
	"os"
	"service/app/queries"

	"service/pkg/configs"
	"service/pkg/logging"
	"service/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var storage queries.Storage


func Init(router *fiber.App, cfg *configs.Config, log *slog.Logger){
	db, err := database.NewClient(context.Background(), cfg)
	if err != nil {
		log.Error("failed to init storage", logging.Err(err))
		os.Exit(0)
	}
	storage = database.NewStorage(db, cfg.MongoDB.Collection)


	serviceGroup := router.Group("/token")
	router.Use(logger.New())
	serviceGroup.Get("/getToken", GetToken)
	serviceGroup.Post("/refreshToken", RefreshToken)

}