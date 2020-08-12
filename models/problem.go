package models

import "github.com/lib/pq"

type Problem struct {
	ID uint `json:"id,omitempty"`

	Name       string         `json:"name,omitempty"`
	Credits    string         `json:"credits,omitempty"`
	Body       string         `json:"body,omitempty"`
	Template   string         `json:"template,omitempty"`
	Difficulty string         `json:"difficulty,omitempty"`
	Tags       pq.StringArray `json:"tags,omitempty"`
}
