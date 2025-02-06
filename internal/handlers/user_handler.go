package handlers

import (
	"net/http"
	"strconv"

	"golang-api-template/internal/i18n"
	"golang-api-template/internal/service"
	"golang-api-template/internal/utils"
	"golang-api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  service.UserService
	emailService *service.EmailService
}

func NewUserHandler(us service.UserService, es *service.EmailService) *UserHandler {
	return &UserHandler{
		userService:  us,
		emailService: es, // Initialize the EmailService here
	}
}

// CREATE
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
		// e.g. userService might return "email is already taken"
		// we can map that if we want, or just show error directly
		// For demonstration, let's do a direct approach:
		if err.Error() == "email is already taken" {
			response.Error(c, http.StatusConflict, i18n.T(c, "EmailTaken"))
			return
		}
		// fallback
		response.Error(c, http.StatusInternalServerError, i18n.T(c, "CreateUserError"))
		return
	}

	response.Success(c, http.StatusCreated, i18n.T(c, "UserCreated"), user)
}

// READ (Get user by ID)
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, i18n.T(c, "InvalidUserID"))
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, i18n.T(c, "UserNotFound"))
		return
	}

	response.Success(c, http.StatusOK, i18n.T(c, "UserRetrieved"), user)
}

// READ (List users)
func (h *UserHandler) List(c *gin.Context) {
	pagination := utils.ParsePagination(c)
	users, total, err := h.userService.GetAllUsers(pagination)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	data := gin.H{
		"pagination": gin.H{
			"page":  pagination.Page,
			"limit": pagination.Limit,
			"total": total,
		},
		"items": users,
	}
	response.Success(c, http.StatusOK, i18n.T(c, "ListOfUsers"), data)
}

// UPDATE
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, i18n.T(c, "InvalidUserID"))
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
		// Check for "email is already taken"
		if err.Error() == "email is already taken by another user" {
			response.Error(c, http.StatusConflict, i18n.T(c, "EmailTaken"))
			return
		}
		response.Error(c, http.StatusInternalServerError, i18n.T(c, "UpdateUserError"))
		return
	}

	response.Success(c, http.StatusOK, i18n.T(c, "UserUpdated"), user)
}

// DELETE
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, i18n.T(c, "InvalidUserID"))
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, i18n.T(c, "DeleteUserError"))
		return
	}

	response.Success(c, http.StatusOK, i18n.T(c, "UserDeleted"), nil)
}

func (h *UserHandler) GetPermissionsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	permissions, err := h.userService.GetPermissionsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(c, "InvalidInput")})
		return
	}

	user, err := h.userService.FindByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": i18n.T(c, "PasswordResetRequestSuccess")})
		return
	}

	resetToken, err := h.userService.GeneratePasswordResetToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(c, "PasswordResetRequestError")})
		return
	}
	resetLink := "https://example.com/reset_password?token=" + resetToken
	lang := c.GetString("locale") // Ensure that the locale middleware sets this.

	go h.emailService.SendPasswordResetEmail(user.Email, resetLink, lang)

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(c, "PasswordResetRequestSuccess")})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(c, "InvalidInput")})
		return
	}

	err := h.userService.ResetPassword(request.Token, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(c, "InvalidOrExpiredToken")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n.T(c, "PasswordResetSuccess")})
}
