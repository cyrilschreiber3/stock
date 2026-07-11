package routes

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/cyrilschreiber3/stock/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	ginLogger := logger.GinLoggerWithConfig("/api/health")
	ginMode := gin.ReleaseMode

	if os.Getenv("ENV") == "localdev" {
		ginLogger = gin.Logger()
	}
	if os.Getenv("LOGLEVEL") == "debug" {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)

	router := gin.New()

	router.Use(ginLogger, gin.Recovery())

	setupCors(router)

	err := router.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatal("Error setting trusted proxies", "error", err)
	}
	return router
}

func setupCors(router *gin.Engine) {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			slog.Debug("Allowed Origin", "origin", origins[i])
		}
	} else {
		origins = []string{"http://localhost:8080", "http://localhost:8081"}
		slog.Debug("Allowed Origin", "origins", origins)
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
}
