package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"` // Use varchar instead of text
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}
