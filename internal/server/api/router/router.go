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
	app.Get("/api/youtube/banner", handler.GetBannerVideos)
	app.Get("/api/youtube/categories", handler.GetAllCategories)
	app.Get("/api/youtube/videos/category/:category", handler.GetVideosByCategory)
	app.Delete("/api/youtube/videos/:video_id/category", handler.RemoveCategoryByID)

	app.Post("/api/youtube/video", handler.InsertOrUpdateVideo)
	app.Get("/api/youtube/video/:videoId", handler.GetVideo)
	app.Delete("/api/youtube/video/:videoId", handler.DeleteVideo)
	app.Post("/api/youtube/search", handler.SearchVideos)

	// Creator Routes
	app.Post("/api/youtube/creator", handler.InsertOrUpdateCreator)
	app.Get("/api/youtube/creator/:creatorId", handler.GetCreator)
	app.Get("/api/youtube/creators", handler.GetCreators)
	app.Get("/api/youtube/creator/:creatorId/videos", handler.GetVideosByCreator)

	app.Get("/s", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"Hello": "world"})
	})

	return app
}
