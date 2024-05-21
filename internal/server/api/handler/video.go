package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	db "github.com/shaik80/ODIW/internal/db/opensearch/controller"
	"github.com/shaik80/ODIW/internal/models"

	"github.com/gofiber/fiber/v2"
)

func InsertOrUpdateVideo(c *fiber.Ctx) error {
	// Parse request body
	var requestBody struct {
		VideoID    string   `json:"video_id"`
		Categories []string `json:"categories"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error parsing request body"})
	}

	videoID := requestBody.VideoID
	if videoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video_id parameter is required"})
	}

	// Send request to fetch video data
	body, err := NewRequest(videoID)
	if err != nil {
		return err
	}

	// Parse JSON response
	var response models.VideoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	// Validate video data
	if err := ValidateVideo(&response.Data); err != nil {
		return err
	}
	response.Data.VideoID = videoID

	// Add categories to the video data
	response.Data.Categories = requestBody.Categories

	// Check if the video already exists in OpenSearch
	existingVideo, _ := db.GetVideoByID(response.Data.VideoID)

	// Insert or update the video
	if existingVideo == nil {
		// Video does not exist, insert it
		if err := db.InsertVideo(&response.Data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		// Video exists, update it if necessary
		isChanged, updatedResponse := CompareAndUpdate(existingVideo, &response.Data)
		if isChanged {
			if err := db.UpdateVideo(updatedResponse); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Video inserted/updated successfully"})
}

// GetVideo retrieves a video by its ID from OpenSearch and returns it in the response
func GetVideo(c *fiber.Ctx) error {
	// Retrieve the video_id parameter from the request
	videoID := c.Params("videoId")
	if videoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video_id parameter is required"})
	}

	// Fetch the video details from OpenSearch using the GetVideoByID function
	video, err := db.GetVideoByID(videoID)
	if err != nil {
		// Check if the error is due to the video not being found
		if err.Error() == "video with ID "+videoID+" not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "youtube video not found"})
		}
		// For other errors, return a 500 Internal Server Error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "youtube video not found"})
	}

	// Return the video details in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"video": video})
}

// SearchVideos searches for videos in the OpenSearch index based on a query parameter with pagination

func SearchVideos(c *fiber.Ctx) error {
	// Parse request body
	var req models.SearchVideosRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate the query parameter
	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter is required"})
	}

	// Set default pagination parameters if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Calculate the starting point for pagination
	from := (req.Page - 1) * req.Size

	// Perform the search operation in the database
	total, videos, err := db.SearchVideos(req.Query, from, req.Size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error searching for videos"})
	}

	// Return the search results with pagination information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":   req.Page,
		"size":   req.Size,
		"total":  total,
		"videos": videos,
	})
}

func GetBannerVideos(c *fiber.Ctx) error {
	_, videos, err := db.SearchVideosByCategory("banner", 0, 10)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "requested data not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"videos": videos,
	})
}

func GetAllCategories(c *fiber.Ctx) error {
	categories, err := db.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"categories": categories})
}

// GetVideosByCategory retrieves videos by a specific category
func GetVideosByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category parameter is required"})
	}

	// Get pagination parameters
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	// Calculate the starting point for pagination
	from := (page - 1) * size

	total, videos, err := db.SearchVideosByCategory(category, from, size)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "requested data not found"})
	}

	// Return the search results with pagination information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":   page,
		"size":   size,
		"total":  total,
		"videos": videos,
	})
}

// RemoveCategoryByID removes a specific category from a video by its ID
func RemoveCategoryByID(c *fiber.Ctx) error {
	videoID := c.Params("video_id")
	if videoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video_id parameter is required"})
	}

	category := c.Query("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "category parameter is required"})
	}

	// Fetch the video by ID
	existingVideo, err := db.GetVideoByID(videoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching video"})
	}

	// Remove the category from the video
	updatedCategories := []string{}
	for _, cat := range existingVideo.Categories {
		if cat != category {
			updatedCategories = append(updatedCategories, cat)
		}
	}
	existingVideo.Categories = updatedCategories

	if err := db.UpdateVideo(existingVideo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Category removed successfully"})
}

// NewRequest sends a GET request to the video info endpoint and returns the response body.
func NewRequest(videoID string) ([]byte, error) {
	url := fmt.Sprintf("https://yig-video-downloader-backend.vercel.app/get_youtube_video_info?url=https://www.youtube.com/watch?v=%s&details=true", videoID)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ValidateVideo validates the video data
func ValidateVideo(video *models.Video) error {
	// Check if video title is empty
	if video.Title == "" {
		return errors.New("title is required")
	}

	// Additional validation rules can be added here

	return nil
}

// CompareAndUpdate checks if the fields of the current video are different from the provided video
// and updates the current video with the new values if necessary.
func CompareAndUpdate(oldVideo, newVideo *models.Video) (bool, *models.Video) {
	changed := false

	if oldVideo.VideoID != newVideo.VideoID {
		fmt.Println("Video ID has changed, updating...")
		oldVideo.VideoID = newVideo.VideoID
		changed = true
	}
	if oldVideo.Title != newVideo.Title {
		fmt.Println("Title has changed, updating...")
		oldVideo.Title = newVideo.Title
		changed = true
	}
	// Compare other fields similarly
	if !compareThumbnails(oldVideo.Thumbnails, newVideo.Thumbnails) {
		fmt.Println("Thumbnails have changed, updating...")
		oldVideo.Thumbnails = newVideo.Thumbnails
		changed = true
	}
	if oldVideo.Likes != newVideo.Likes {
		fmt.Println("Likes count has changed, updating...")
		oldVideo.Likes = newVideo.Likes
		changed = true
	}
	if oldVideo.ViewsCount != newVideo.ViewsCount {
		fmt.Println("Views count has changed, updating...")
		oldVideo.ViewsCount = newVideo.ViewsCount
		changed = true
	}
	// Add comparisons for other fields as needed
	if oldVideo.UploadDate != newVideo.UploadDate {
		fmt.Println("Upload date has changed, updating...")
		oldVideo.UploadDate = newVideo.UploadDate
		changed = true
	}
	if oldVideo.VideoCategory != newVideo.VideoCategory {
		fmt.Println("Video category has changed, updating...")
		oldVideo.VideoCategory = newVideo.VideoCategory
		changed = true
	}
	if oldVideo.Description != newVideo.Description {
		fmt.Println("Description has changed, updating...")
		oldVideo.Description = newVideo.Description
		changed = true
	}
	if oldVideo.Dislikes != newVideo.Dislikes {
		fmt.Println("Dislikes count has changed, updating...")
		oldVideo.Dislikes = newVideo.Dislikes
		changed = true
	}
	if oldVideo.IsShort != newVideo.IsShort {
		fmt.Println("IsShort flag has changed, updating...")
		oldVideo.IsShort = newVideo.IsShort
		changed = true
	}
	if oldVideo.CreatorDetails != newVideo.CreatorDetails {
		fmt.Println("Creator details have changed, updating...")
		oldVideo.CreatorDetails = newVideo.CreatorDetails
		changed = true
	}
	if oldVideo.LastUpdated != newVideo.LastUpdated {
		fmt.Println("Last updated timestamp has changed, updating...")
		oldVideo.LastUpdated = newVideo.LastUpdated
		changed = true
	}

	return changed, oldVideo
}

// compareThumbnails compares two slices of Thumbnail and returns true if they are equal, false otherwise.
func compareThumbnails(a, b []models.Thumbnail) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
