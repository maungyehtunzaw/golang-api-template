package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"golang-api-template/internal/config"
	"golang-api-template/internal/models"
	"golang-api-template/internal/repository"
	"golang-api-template/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(name, email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers(p utils.PaginationParams) ([]models.User, int64, error)
	UpdateUser(id uint, name, email, password string) (*models.User, error)
	DeleteUser(id uint) error
	GetPermissionsByUserID(userID uint) ([]models.Permission, error)

	FindByEmail(email string) (*models.User, error)
	GeneratePasswordResetToken(user *models.User) (string, error)
	ResetPassword(token, newPassword string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		repo: r,
	}
}

// CREATE
func (s *userService) CreateUser(name, email, password string) (*models.User, error) {
	// Check if user already exists by email
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("email is already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// READ (single)
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

// READ (all with pagination)
func (s *userService) GetAllUsers(p utils.PaginationParams) ([]models.User, int64, error) {
	return s.repo.GetUsers(p)
}

// UPDATE
func (s *userService) UpdateUser(id uint, name, email, password string) (*models.User, error) {
	// 1. Find the user
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Update fields
	user.Name = name
	// Check if new email is already taken by another user
	if user.Email != email {
		if _, err := s.repo.GetUserByEmail(email); err == nil {
			return nil, errors.New("email is already taken by another user")
		}
		user.Email = email
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// 3. Save updated user
	err = s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DELETE
func (s *userService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) GetPermissionsByUserID(userID uint) ([]models.Permission, error) {
	return s.repo.GetPermissionsByUserID(userID)
}

func (s *userService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *userService) GeneratePasswordResetToken(user *models.User) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err // handle random generator error
	}
	token := fmt.Sprintf("%x", b)
	expiry := time.Now().Add(config.GetResetTokenExpiry()) // Assuming config package has this method

	// Store the token in the database with an expiration time
	if err := s.repo.SaveResetToken(user.ID, token, expiry); err != nil {
		return "", err // handle database error
	}

	return token, nil
}
func (s *userService) ResetPassword(token, newPassword string) error {
	user, err := s.repo.FindByToken(token)
	if err != nil {
		return err
	}

	// Additional checks, such as token expiration, can be implemented here
	hashedPassword, _ := utils.HashPassword(newPassword) // Assuming you have a HashPassword method
	return s.repo.UpdatePassword(user.ID, hashedPassword)
}
