package repository

import (
	"golang-api-template/internal/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *models.Role) error
	GetAllRoles() ([]models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error
	GetPermissionsByRoleID(roleID uint) ([]models.Permission, error)
}

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepo{db}
}

func (repo *roleRepo) CreateRole(role *models.Role) error {
	return repo.db.Create(role).Error
}

func (repo *roleRepo) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := repo.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (repo *roleRepo) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := repo.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (repo *roleRepo) UpdateRole(role *models.Role) error {
	return repo.db.Save(role).Error
}

func (repo *roleRepo) DeleteRole(id uint) error {
	return repo.db.Delete(&models.Role{}, id).Error
}
func (repo *roleRepo) GetPermissionsByRoleID(roleID uint) ([]models.Permission, error) {
	var role models.Role
	if err := repo.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
