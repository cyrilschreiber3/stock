package main

import (
	"log/slog"

	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/cyrilschreiber3/stock/logger"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/cyrilschreiber3/stock/utils"
)

func main() {
	utils.LoadEnv()
	defaultLogger := logger.InitLogger()
	slog.SetDefault(defaultLogger)
	slog.Info("Stock - Stock management application")

	cmd := utils.GetCommand()
	switch cmd {
	case "migrate":
		database.Migrate()
		return
	case "serve":
		slog.Info("Starting server")
	default:
		logger.Fatal("Unknown command", "command", cmd)
	}
	database.Init()
	defer database.Close()

	handlers.Init()

	router := routes.Init()
	routes.RegisterRoutes(router)

	slog.Info("Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
