package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"size:100" json:"name"`              // Limit Name to 100 characters
	Email       string     `gorm:"size:255;uniqueIndex" json:"email"` // Limit Email to 255 characters
	Password    string     `gorm:"size:225" json:"-"`
	Roles       []Role     `gorm:"many2many:user_roles" json:"roles"`              // Hashed password, omit from JSON
	ResetToken  string     `gorm:"index"`                                          // Index this field for faster lookup
	TokenExpiry *time.Time `gorm:"type:datetime;default:null" json:"token_expiry"` // Make TokenExpiry nullable

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"` // Auto-set at creation
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"` // Auto-set on update
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // Soft delete support
}
