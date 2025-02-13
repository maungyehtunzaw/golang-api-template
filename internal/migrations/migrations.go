package migrations

import (
	"golang-api-template/internal/models"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM's automigration for all models
func AutoMigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
}
