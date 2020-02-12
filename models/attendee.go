package models

import (
	"github.com/jinzhu/gorm"
)

type Attendee struct {
	gorm.Model
	Name    string
	EventID uint
}
