package model

import "time"

// Package represents a knowledge transfer package.
type Package struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary,omitempty"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	Project   string    `json:"project,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
