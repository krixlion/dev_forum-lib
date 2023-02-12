package testtypes

import "time"

type Article struct {
	Id        string    `json:"id,omitempty"`
	UserId    string    `json:"user_id,omitempty"` // Author's ID.
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
