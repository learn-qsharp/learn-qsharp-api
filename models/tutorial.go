package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Tutorial struct {
	gorm.Model
	Title       string         `valid:"required"`
	Author      string         `valid:"required"`
	Description string         `valid:"required"`
	Difficulty  string         `valid:"required,in(easy|medium|hard)"`
	Tags        pq.StringArray `gorm:"type:varchar(64)[]"`
}
