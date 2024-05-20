package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	connect "shaik80/ODIW/internal/db/opensearch"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// InsertDocument inserts a document into the specified index
func InsertDocument(index, documentID string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling document: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return fmt.Errorf("error inserting document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return fmt.Errorf("failed to insert document, status: %s", res.Status())
	}

	return nil
}

// UpdateDocument updates a document in the specified index
func UpdateDocument(index, documentID string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling document: %w", err)
	}

	req := opensearchapi.UpdateRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to update document, status: %s", res.Status())
	}

	return nil
}

// DeleteDocument deletes a document from the specified index
func DeleteDocument(index, documentID string) error {
	req := opensearchapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to delete document, status: %s", res.Status())
	}

	return nil
}

// SearchDocuments searches for documents in the specified index with the given query
func SearchDocuments(index, query string) (*opensearchapi.Response, error) {
	req := opensearchapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), connect.Client)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %w", err)
	}

	return res, nil
}
