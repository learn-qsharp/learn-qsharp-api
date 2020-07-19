package models

import (
	"github.com/lib/pq"
)

type Tutorial struct {
	ID uint `json:"id,omitempty"`

	Title       string         `json:"title,omitempty"`
	Credits     string         `json:"credits,omitempty"`
	Description string         `json:"description,omitempty"`
	Body        string         `json:"body,omitempty" `
	Difficulty  string         `json:"difficulty,omitempty" `
	Tags        pq.StringArray `json:"tags,omitempty"`
}
