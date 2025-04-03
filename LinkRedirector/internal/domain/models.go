package domain

import "time"

type UrlDAO struct {
	ID        int       `json:"id"`
	Original  string    `json:"original"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
