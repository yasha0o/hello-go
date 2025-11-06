package dto

import (
	"time"

	"hello-go/internal/domain"

	"github.com/google/uuid"
)

type DocumentDto struct {
	ID        uuid.UUID `json:"id"`
	Document  any       `json:"document"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02T15:04:05Z"`
}

type PageParams struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type RequestedDocumentsPage struct {
	Data  []DocumentDto `json:"data"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

type DocumentRequest struct {
	DocumentID string        `uri:"id" binding:"required"`
	Host       string        `header:"Host"`
	Sender     string        `form:"sender" binding:"required"`
	Region     domain.Region `form:"region"`
}
