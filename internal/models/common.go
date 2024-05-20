package models

type SearchVideosRequest struct {
	Query string `json:"query" validate:"required"`
	Page  int    `json:"page" validate:"required,min=1"`
	Size  int    `json:"size" validate:"required,min=1"`
}
