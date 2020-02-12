package models

import (
	"github.com/jinzhu/gorm"
)

type Guest struct {
	gorm.Model
	Name           *string `gorm:"not null"`
	EventDisplayID *string `gorm:"not null"`
}
