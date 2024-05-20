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

	videoID := c.Query("video_id")
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

	// Check if the video already exists in OpenSearch
	existingVideo, err := db.GetVideoByID(response.Data.VideoID)
	if err != nil {
		return err
	}
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

func GetVideo(c *fiber.Ctx) error {
	return nil
}

func SearchVideos(c *fiber.Ctx) error {
	return nil
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
