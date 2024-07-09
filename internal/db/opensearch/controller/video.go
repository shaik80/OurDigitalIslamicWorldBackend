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

// SearchVideos queries the OpenSearch index for videos matching the query with pagination
func SearchVideos(query string, from int, size int) (int, []*models.Video, error) {
	// Create search request with pagination
	searchRequest := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "description", "tags", "categories"},
			},
		},
		"track_total_hits": true, // Ensure total hits is tracked
	}

	// Serialize search request to JSON
	searchData, err := json.Marshal(searchRequest)
	if err != nil {
		return 0, nil, err
	}

	// Perform search request
	searchReq := strings.NewReader(string(searchData))
	// Perform search request
	res, err := connect.Client.Search(
		context.Background(),
		&opensearchapi.SearchReq{
			Indices: []string{"videos"},
			Body:    searchReq,
		},
	)
	if err != nil {
		return 0, nil, err
	}
	fmt.Printf("Search hits: %v\n", res.Hits.Total.Value)

	// Check if the search response is an error
	if res.Inspect().Response.IsError() {
		return 0, nil, fmt.Errorf("search response error: %v", res.Errors)
	}

	// Parse search response
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Inspect().Response.Body).Decode(&searchResponse); err != nil {
		return 0, nil, err
	}

	// Extract video results from search response
	hits := searchResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	videos := make([]*models.Video, len(hits))
	for i, hit := range hits {
		videoData := hit.(map[string]interface{})["_source"]
		videoBytes, _ := json.Marshal(videoData)
		var video models.Video
		if err := json.Unmarshal(videoBytes, &video); err != nil {
			return 0, nil, err
		}
		videos[i] = &video
	}
	// total := searchResponse["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(int),
	// Extract total hits count
	totalHits := searchResponse["hits"].(map[string]interface{})["total"]
	var total int
	switch totalHits := totalHits.(type) {
	case map[string]interface{}:
		total = int(totalHits["value"].(float64))
	case float64:
		total = int(totalHits)
	}
	return total, videos, nil
}

func GetAllCategories() ([]string, error) {
	// Create a search request to get all categories
	searchRequest := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"unique_categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "categories.keyword",
					"size":  1000,
				},
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
			Indices: []string{"videos"},
			Body:    strings.NewReader(string(searchData)),
		},
	)
	if err != nil {
		return nil, err
	}

	// Check if the search response is an error
	if res.Inspect().Response.IsError() {
		return nil, fmt.Errorf("search response error: %v", res.Errors)
	}

	// Parse search response
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Inspect().Response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	// Extract unique categories from search response
	buckets := searchResponse["aggregations"].(map[string]interface{})["unique_categories"].(map[string]interface{})["buckets"].([]interface{})
	categories := make([]string, len(buckets))
	for i, bucket := range buckets {
		categories[i] = bucket.(map[string]interface{})["key"].(string)
	}

	return categories, nil
}

// SearchVideosByCategory queries the OpenSearch index for videos matching the category with pagination
func SearchVideosByCategory(category string, from int, size int) (int, []*models.Video, error) {
	// Create search request with pagination
	searchRequest := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"categories": category,
			},
		},
		"track_total_hits": true, // Ensure total hits is tracked
	}

	// Serialize search request to JSON
	searchData, err := json.Marshal(searchRequest)
	if err != nil {
		return 0, nil, err
	}

	// Perform search request
	searchReq := strings.NewReader(string(searchData))

	// Perform search request
	res, err := connect.Client.Search(
		context.Background(),
		&opensearchapi.SearchReq{
			Indices: []string{"videos"},
			Body:    searchReq,
		},
	)
	if err != nil {
		return 0, nil, err
	}

	// Check if the search response is an error
	if res.Inspect().Response.IsError() {
		return 0, nil, fmt.Errorf("search response error: %v", res.Errors)
	}

	// Parse search response
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Inspect().Response.Body).Decode(&searchResponse); err != nil {
		return 0, nil, err
	}

	// Extract video results from search response
	hits := searchResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	videos := make([]*models.Video, len(hits))
	for i, hit := range hits {
		videoData := hit.(map[string]interface{})["_source"]
		videoBytes, _ := json.Marshal(videoData)
		var video models.Video
		if err := json.Unmarshal(videoBytes, &video); err != nil {
			return 0, nil, err
		}
		videos[i] = &video
	}

	// Extract total hits count
	totalHits := searchResponse["hits"].(map[string]interface{})["total"]
	var total int
	switch totalHits := totalHits.(type) {
	case map[string]interface{}:
		total = int(totalHits["value"].(float64))
	case float64:
		total = int(totalHits)
	}

	return total, videos, nil
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
	ctx := context.Background()
	// Check if the "videos" index exists in the OpenSearch cluster
	indexExists, err := connect.Client.Indices.Exists(ctx, opensearchapi.IndicesExistsReq{
		Indices: []string{"videos"},
	})
	if err != nil && indexExists.StatusCode != 404 {
		return nil, fmt.Errorf("error checking if index exists: %w", err)
	}

	if indexExists.StatusCode != 200 {
		_, err := connect.Client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
			Index: "videos",
			Body:  strings.NewReader(""),
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating index")
		}
	}

	// Create a request to retrieve the document by ID
	req := opensearchapi.DocumentGetReq{
		Index:      "videos",
		DocumentID: videoID,
	}

	getResponse, err := connect.Client.Document.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// Check if the document exists
	if getResponse.Inspect().Response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("video with ID %s not found", videoID)
	} else if getResponse.Inspect().Response.IsError() {
		return nil, fmt.Errorf("error getting video with ID %s: %s", videoID, getResponse.Inspect().Response.Status())
	}

	// Parse the response body into a video object
	var video models.Video
	err = json.Unmarshal(getResponse.Source, &video)
	if err != nil {
		return nil, fmt.Errorf("error decoding video response: %s", err)
	}

	return &video, nil
}

func DeleteVideoByID(videoID string) error {
	// Create delete request
	req := opensearchapi.DocumentDeleteReq{
		Index:      "videos",
		DocumentID: videoID,
		Params: opensearchapi.DocumentDeleteParams{
			Refresh: "true",
		},
	}

	deleteResponse, err := connect.Client.Document.Delete(context.Background(), req)
	if err != nil {
		return err
	}

	// Check if the delete operation was successful
	if deleteResponse.Inspect().Response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("video with ID %s not found", videoID)
	} else if deleteResponse.Inspect().Response.IsError() {
		return fmt.Errorf("error deleting video with ID %s: %s", videoID, deleteResponse.Inspect().Response.Status())
	}

	fmt.Printf("Deleted document with ID %s\n", videoID)
	return nil
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
