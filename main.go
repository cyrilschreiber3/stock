package main

import (
	"log/slog"

	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/cyrilschreiber3/stock/logger"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()
	logger := logger.InitLogger()
	slog.SetDefault(logger)
	slog.Info("Starting stock application")

	database.Init()
	defer database.Close()

	// initialize handlers which wires the sqlc-generated queries instance
	handlers.Init()

	router := gin.Default()
	utils.SetupRouter(router)
	routes.SetupRoutes(router)

	slog.Info("Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
