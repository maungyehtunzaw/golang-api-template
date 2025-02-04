package handlers

import (
	"net/http"
	"strconv"

	"golang-api-template/internal/service"
	"golang-api-template/internal/utils"
	"golang-api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// ------------------------------------
// CREATE (Register a new user)
// ------------------------------------
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error()) // e.g., email taken
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

// ------------------------------------
// READ (Get user by ID)
// ------------------------------------
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// ------------------------------------
// READ (List users with pagination)
// ------------------------------------
func (h *UserHandler) List(c *gin.Context) {
	pagination := utils.ParsePagination(c)

	users, total, err := h.userService.GetAllUsers(pagination)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Example: You can either do a simple response with
	// "pagination" and "items" or use a "Laravel-style" approach.
	// For simplicity, we'll do a simple approach here:
	data := gin.H{
		"pagination": gin.H{
			"page":  pagination.Page,
			"limit": pagination.Limit,
			"total": total,
		},
		"items": users,
	}

	response.Success(c, http.StatusOK, "List of users", data)
}

// ------------------------------------
// UPDATE (Update user by ID)
// ------------------------------------
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// ------------------------------------
// DELETE (Delete user by ID)
// ------------------------------------
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
