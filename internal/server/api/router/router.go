package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shaik80/ODIW/internal/server/api/handler"
)

// Setup initializes and returns the Fiber app with defined routes
func SetupRoutes(app *fiber.App) *fiber.App {
	// Initialize Handler
	// register Handlers

	// routers

	app.Post("/api/youtube/video", handler.InsertOrUpdateVideo)
	app.Get("/api/youtube/video/:videoId", handler.GetVideo)
	app.Post("/api/youtube/search", handler.SearchVideos)

	// Creator Routes
	app.Post("/api/youtube/creator", handler.InsertOrUpdateCreator)
	app.Get("/api/youtube/creator/:creatorId", handler.GetCreator)
	app.Get("/api/youtube/creators", handler.GetCreators)
	app.Get("/api/youtube/creator/:creatorId/videos", handler.GetVideosByCreator)

	return app
}
