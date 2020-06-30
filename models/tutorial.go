package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Tutorial struct {
	gorm.Model
	Title       string
	Author      string
	Description string
	Difficulty  string
	Tags        pq.StringArray `gorm:"type:varchar(64)[]"`
}
