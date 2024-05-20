package router

import (
	"shaik80/ODIW/internal/server/api/handler"

	"github.com/gofiber/fiber/v2"
)

// Setup initializes and returns the Fiber app with defined routes
func Setup(app *fiber.App) *fiber.App {
	// Initialize Handler
	// register Handlers

	// routers

	app.Get("/api/youtube/video", handler.InsertOrUpdateVideo)
	app.Get("/api/youtube/video/:videoId", handler.GetVideo)
	app.Get("/api/youtube/search", handler.SearchVideos)

	// Creator Routes
	app.Post("/api/youtube/creator", handler.InsertOrUpdateCreator)
	app.Get("/api/youtube/creator/:creatorId", handler.GetCreator)
	app.Get("/api/youtube/creators", handler.GetCreators)
	app.Get("/api/youtube/creator/:creatorId/videos", handler.GetVideosByCreator)

	return app
}
