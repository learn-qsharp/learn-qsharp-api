package models

import (
	"github.com/jinzhu/gorm"
)

type Tutorial struct {
	gorm.Model
	Title       string
	Description string
	Difficulty  int
}
