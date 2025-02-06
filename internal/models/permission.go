package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name  string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"` // Use varchar instead of text
	Roles []Role `gorm:"many2many:role_permissions" json:"roles"`
}
