package models

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CollectionRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *CollectionRequest) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type CollectionBook struct {
	CollectionID int `json:"collection_id"`
	BookID int `json:"book_id"`
}
