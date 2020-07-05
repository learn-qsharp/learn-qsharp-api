package models

import (
	"github.com/lib/pq"
	"time"
)

type Tutorial struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`

	Title       string         `json:"title" valid:"required"`
	Author      string         `json:"author" valid:"required"`
	Description string         `json:"description" valid:"required"`
	Difficulty  string         `json:"difficulty" valid:"required,in(easy|medium|hard)"`
	Tags        pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
}
