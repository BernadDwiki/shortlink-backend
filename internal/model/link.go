package model

import "time"

type Link struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	OriginalURL string     `json:"original_url"`
	Slug        string     `json:"slug"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
