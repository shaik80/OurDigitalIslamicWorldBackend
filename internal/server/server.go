package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shaik80/ODIW/config"
	"github.com/shaik80/ODIW/internal/server/api/router"
)

func SetupGofiber() {
	// Create a new Fiber instance
	app := fiber.New()
	useCustomLoggerMiddleware(app)

	// Setup routes
	router.SetupRoutes(app)
	// Start the server on port 3000
	log.Fatal(app.Listen(fmt.Sprint(":", config.Cfg.Server.Port)))
}

func useCustomLoggerMiddleware(app *fiber.App) {
	// Define custom logger format and colors
	loggerConfig := logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05", // Human-readable time format
		TimeZone:   "Local",
		Output:     os.Stdout,
	}

	// Create a custom logger middleware
	customLogger := logger.New(loggerConfig)

	// Use the custom logger middleware
	app.Use(customLogger)
}
