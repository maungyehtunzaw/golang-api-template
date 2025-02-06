package service

import (
	"golang-api-template/internal/models"
	"golang-api-template/internal/repository"
)

type RoleService interface {
	CreateRole(role *models.Role) error
	GetAllRoles() ([]models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error
	GetPermissionsByRoleID(roleID uint) ([]models.Permission, error)
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo}
}

func (s *roleService) CreateRole(role *models.Role) error {
	return s.repo.CreateRole(role)
}

func (s *roleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAllRoles()
}

func (s *roleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.repo.GetRoleByID(id)
}

func (s *roleService) UpdateRole(role *models.Role) error {
	return s.repo.UpdateRole(role)
}

func (s *roleService) DeleteRole(id uint) error {
	return s.repo.DeleteRole(id)
}
func (s *roleService) GetPermissionsByRoleID(roleID uint) ([]models.Permission, error) {
	return s.repo.GetPermissionsByRoleID(roleID)
}
