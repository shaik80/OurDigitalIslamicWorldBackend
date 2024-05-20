package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	connect "shaik80/ODIW/internal/db/opensearch"
	"shaik80/ODIW/internal/models"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func InsertOrUpdateVideo(video *models.Video) error {
	// Serialize video object to JSON
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}

	// Create request to index document in OpenSearch
	req := opensearchapi.IndexRequest{
		Index:      "videos",
		DocumentID: video.VideoID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true", // Refresh index after operation
	}

	// Execute request
	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// Handle error response
		return fmt.Errorf("failed to index document: %s", res.String())
	}

	return nil
}

func SearchVideos(query string) ([]*models.Video, error) {
	// Create search request
	searchRequest := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": query,
			},
		},
	}

	// Serialize search request to JSON
	searchData, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	// Perform search request
	res, err := connect.Client.Search(
		connect.Client.Search.WithIndex("videos"),
		connect.Client.Search.WithBody(strings.NewReader(string(searchData))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse search response
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	// Extract video results from search response
	var videos []*models.Video
	// Implement extraction logic based on your response structure
	// Example: videos := extractVideos(searchResponse)

	return videos, nil
}

func DeleteVideo(videoID string) error {
	// Create delete request
	req := opensearchapi.DeleteRequest{
		Index:      "videos",
		DocumentID: videoID,
		Refresh:    "true", // Refresh index after operation
	}

	// Execute request
	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// Handle error response
		return fmt.Errorf("failed to delete document: %s", res.String())
	}

	return nil
}

func GetVideoByID(videoID string) (*models.Video, error) {
	// Create a request to retrieve the document by ID
	req := opensearchapi.GetRequest{
		Index:      "videos",
		DocumentID: videoID,
	}

	fmt.Println(context.Background(), connect.Client)

	// Execute the request
	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		fmt.Println(err, "err")

		return nil, err
	}
	defer res.Body.Close()

	// Check if the document exists
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("video with ID %s not found", videoID)
	} else if res.IsError() {
		return nil, fmt.Errorf("error getting video with ID %s: %s", videoID, res.Status())
	}

	// Parse the response body into a video object
	var video models.Video
	if err := json.NewDecoder(res.Body).Decode(&video); err != nil {
		return nil, fmt.Errorf("error decoding video response: %s", err)
	}

	return &video, nil
}

func InsertVideo(video *models.Video) error {
	// Serialize video object to JSON
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}

	// Create request to index document in OpenSearch
	req := opensearchapi.IndexRequest{
		Index:      "videos",
		DocumentID: video.VideoID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true", // Refresh index after operation
	}

	// Execute request
	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// Handle error response
		return fmt.Errorf("failed to index document: %s", res.String())
	}

	return nil
}

func UpdateVideo(video *models.Video) error {
	// Serialize video object to JSON
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}

	// Create request to update document in OpenSearch
	req := opensearchapi.UpdateRequest{
		Index:      "videos",
		DocumentID: video.VideoID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true", // Refresh index after operation
	}

	// Execute request
	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// Handle error response
		return fmt.Errorf("failed to update document: %s", res.String())
	}

	return nil
}
