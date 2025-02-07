package handlers

import (
	"fmt"
	"net/http"

	"golang-api-template/internal/i18n"
	"golang-api-template/internal/service"
	"golang-api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use our Error response with a 400 status code
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		// Use our Error response with a 401 status code
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.authService.TrackUserLogin(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Errorf("failed to track user login: %w", err).Error())
		return
	}

	// Use our Success response with 200 status code
	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	newAccess, newRefresh, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) GetAuthUser(c *gin.Context) {
	// you can use here to get user
	userId, exists := c.Get("AuthID")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Error getting user ID")
		return
	}
	fmt.Println("User ID:", userId)
	user, err := h.authService.GetAuthUser(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	userResponse := response.AuthUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	response.Success(c, http.StatusOK, i18n.T(c, "UserFound"), userResponse)
}
