package models

import (
	"fmt"
	"time"
)

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate string    `json:"published_date"`
	Edition       int       `json:"edition"`
	Description   string    `json:"description"`
	Genre         string    `json:"genre"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookRequest struct {
	Title         string `json:"title,omitempty"`
	Author        string `json:"author,omitempty"`
	PublishedDate string `json:"published_date,omitempty"`
	Edition       int    `json:"edition,omitempty"`
	Description   string `json:"description,omitempty"`
	Genre         string `json:"genre,omitempty"`
}

type BookGroup struct {
	GroupKey string `json:"group_key"`
	Count    int    `json:"count"`
	Books    []Book `json:"books"`
}

func (b *BookRequest) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("title is required")
	}
	if b.Author == "" {
		return fmt.Errorf("author is required")
	}
	if b.PublishedDate == "" {
		return fmt.Errorf("published_date is required")
	}

	if b.PublishedDate != "" {
		if _, err := time.Parse("2006-01-02", b.PublishedDate); err != nil {
			return fmt.Errorf("invalid published_date format")
		}
	}
	return nil
}
