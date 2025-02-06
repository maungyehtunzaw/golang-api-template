package repository

import (
	"golang-api-template/internal/models"
	"golang-api-template/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error

	// For listing with pagination
	GetUsers(p utils.PaginationParams) ([]models.User, int64, error)
	GetPermissionsByUserID(userID uint) ([]models.Permission, error)

	FindByEmail(email string) (*models.User, error)
	SaveResetToken(userID uint, token string, expiry time.Time) error // Corrected signature

	FindByToken(token string) (*models.User, error)
	UpdatePassword(userID uint, newPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetUsers fetches a list of users with pagination
func (r *userRepository) GetUsers(p utils.PaginationParams) ([]models.User, int64, error) {
	var (
		users []models.User
		total int64
	)

	// Count total
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	if err := r.db.Offset(p.Offset()).Limit(p.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (repo *userRepository) GetPermissionsByUserID(userID uint) ([]models.Permission, error) {
	var user models.User
	if err := repo.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	var permissions []models.Permission
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			// Avoid adding duplicate permissions
			if !containsPermission(permissions, perm) {
				permissions = append(permissions, perm)
			}
		}
	}
	return permissions, nil
}

func containsPermission(permissions []models.Permission, permission models.Permission) bool {
	for _, perm := range permissions {
		if perm.ID == permission.ID {
			return true
		}
	}
	return false
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) SaveResetToken(userID uint, token string, expiry time.Time) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(models.User{ResetToken: token, TokenExpiry: expiry}).Error
}

func (r *userRepository) FindByToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID uint, newPassword string) error {
	// Example hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
