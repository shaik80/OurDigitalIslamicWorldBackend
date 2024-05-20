package server

import (
	"shaik80/ODIW/internal/server/api/router"

	"github.com/gofiber/fiber/v2"
	api_log "github.com/gofiber/fiber/v2/middleware/logger"
)

func New() *fiber.App {

	app := fiber.New()

	// Setup handlers and routes
	// handlers.SetupHandlers(app)
	// logger.Logger.Printf("Server is running on %s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)
	// Define your routes here
	app.Use(api_log.New())

	return router.Setup(app)
}
