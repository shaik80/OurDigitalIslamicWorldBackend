package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	connect "github.com/shaik80/ODIW/internal/db/opensearch"
	"github.com/shaik80/ODIW/internal/models"
)

func InsertOrUpdateVideo(video *models.Video) error {
	// Serialize video object to JSON
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}

	// Create request to index document in OpenSearch
	req := opensearchapi.IndexReq{
		Index:      "videos",
		DocumentID: video.VideoID,
		Body:       strings.NewReader(string(data)),
		Params: opensearchapi.IndexParams{
			Refresh: "true",
		},
		// Refresh:    "true", // Refresh index after operation
	}

	insertResp, err := connect.Client.Index(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

	// Execute request
	// res, err := req.Do(context.Background(), connect.Client)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()

	// if res.IsError() {
	// 	// Handle error response
	// 	return fmt.Errorf("failed to index document: %s", res.String())
	// }

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
		context.Background(),
		&opensearchapi.SearchReq{
			Body: strings.NewReader(string(searchData)),
		},
		// connect.Client.Search.WithIndex("videos"),
		// connect.Client.Search.WithBody(strings.NewReader(string(searchData))),
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Search hits: %v\n", res.Hits.Total.Value)

	// defer res.Body.Close()

	// Parse search response
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Inspect().Response.Body).Decode(&searchResponse); err != nil {
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
	req := opensearchapi.DocumentDeleteReq{
		Index:      "videos",
		DocumentID: videoID,
		Params: opensearchapi.DocumentDeleteParams{
			Refresh: "true",
		}}
	deleteResponse, err := connect.Client.Document.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	// Execute request
	// res, err := req.Do(context.Background(), connect.Client)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()

	// if res.IsError() {
	// 	// Handle error response
	// 	return fmt.Errorf("failed to delete document: %s", res.String())
	// }
	fmt.Printf("Deleted document: %t\n", deleteResponse.Result == "deleted")

	return nil
}

func GetVideoByID(videoID string) (*models.Video, error) {
	// Create a request to retrieve the document by ID
	req := opensearchapi.DocumentGetReq{
		Index:      "videos",
		DocumentID: videoID,
	}

	getResponse, err := connect.Client.Document.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	// Execute the request
	// res, err := req.Do(context.Background(), connect.Client)
	// if err != nil {
	// 	fmt.Println(err, "err")

	// 	return nil, err
	// }
	// defer res.Body.Close()
	fmt.Printf("getresponse document: %t\n", getResponse.Found)

	// Check if the document exists
	if getResponse.Inspect().Response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("video with ID %s not found", videoID)
	} else if getResponse.Inspect().Response.IsError() {
		return nil, fmt.Errorf("error getting video with ID %s: %s", videoID, getResponse.Inspect().Response.Status())
	}

	// Parse the response body into a video object
	var video models.Video
	if err := json.NewDecoder(getResponse.Inspect().Response.Body).Decode(&video); err != nil {
		return nil, fmt.Errorf("error decoding video response: %s", err)
	}

	return &video, nil
}

func InsertVideo(videos *models.Video) error {
	// Serialize video object to JSON
	data, err := json.Marshal(videos)
	if err != nil {
		return err
	}

	// Create request to index document in OpenSearch
	req := opensearchapi.IndexReq{
		Index:      "videos",
		DocumentID: videos.VideoID,
		Body:       strings.NewReader(string(data)),
		Params: opensearchapi.IndexParams{
			Refresh: "true",
		},
	}

	// Execute request
	insertResp, err := connect.Client.Index(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

	if insertResp.Inspect().Response.IsError() {
		// Handle error response
		return fmt.Errorf("failed to update document: %s", insertResp.Inspect().Response.String())
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
	req := opensearchapi.IndexReq{
		Index:      "videos",
		DocumentID: video.VideoID,
		Body:       strings.NewReader(string(data)),
		Params: opensearchapi.IndexParams{
			Refresh: "true",
		},
	}

	// Execute request
	// res, err := req.Do(context.Background(), connect.Client)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()
	updateResp, err := connect.Client.Index(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("Created document in %s\n  ID: %s\n", updateResp.Index, updateResp.ID)

	if updateResp.Inspect().Response.IsError() {
		// Handle error response
		return fmt.Errorf("failed to update document: %s", updateResp.Inspect().Response.String())
	}

	return nil
}
