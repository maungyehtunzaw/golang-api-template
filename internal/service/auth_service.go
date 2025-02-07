package service

import (
	"context"
	"fmt"
	"time"

	"golang-api-template/internal/config"
	"golang-api-template/internal/models"
	"golang-api-template/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, string, *models.User, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(refreshToken string) error
	GetAuthUser(ctx context.Context) (*models.User, error)
	TrackUserLogin(userID uint) error
	TrackUserLogout(userID uint) error
	IsUserOnline(userID uint) (bool, error)
}

type authService struct {
	userRepo repository.UserRepository
	rdb      *redis.Client
	cfg      *config.Config
}

func NewAuthService(repo repository.UserRepository, rdb *redis.Client, cfg *config.Config) AuthService {
	return &authService{
		userRepo: repo,
		rdb:      rdb,
		cfg:      cfg,
	}
}

// ----------------------------------------------------------
// LOGIN
// ----------------------------------------------------------
func (s *authService) Login(email, password string) (string, string, *models.User, error) {
	// 1. Find user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid credentials")
	}

	// 2. Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, fmt.Errorf("invalid credentials")
	}

	// 3. Create access token
	accessToken, err := s.createToken(user.ID, s.cfg.JWTAccessSecret, time.Minute*time.Duration(s.cfg.AccessTokenExpireMin))
	if err != nil {
		return "", "", nil, err
	}

	// 4. Create refresh token
	refreshToken, err := s.createToken(user.ID, s.cfg.JWTRefreshSecret, time.Hour*time.Duration(s.cfg.RefreshTokenExpireHrs))
	if err != nil {
		return "", "", nil, err
	}

	// 5. Store the refresh token in Redis
	ctx := context.Background()
	// We store the token as key -> userID, so we can reference it later
	err = s.rdb.Set(ctx, refreshToken, user.ID, time.Hour*time.Duration(s.cfg.RefreshTokenExpireHrs)).Err()
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

// ----------------------------------------------------------
// REFRESH TOKEN
// ----------------------------------------------------------
func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// 1. Validate refresh token signature
	claims, err := s.validateToken(refreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	// 2. Check if refresh token is in Redis
	ctx := context.Background()
	val, err := s.rdb.Get(ctx, refreshToken).Result()
	if err == redis.Nil || val == "" {
		return "", "", fmt.Errorf("refresh token not found or expired")
	} else if err != nil {
		return "", "", err
	}

	// 3. Create new access token
	accessToken, err := s.createToken(uint(userID), s.cfg.JWTAccessSecret, time.Minute*time.Duration(s.cfg.AccessTokenExpireMin))
	if err != nil {
		return "", "", err
	}

	// You can decide whether to issue a new refresh token or keep using the old one
	// Typically, you might re-issue a new refresh token and delete the old one.
	// For simplicity, let's keep the same refresh token.

	return accessToken, refreshToken, nil
}

// ----------------------------------------------------------
// LOGOUT
// ----------------------------------------------------------
func (s *authService) Logout(refreshToken string) error {
	// remove the refresh token from Redis so it can’t be used
	ctx := context.Background()
	err := s.rdb.Del(ctx, refreshToken).Err()
	return err
}

// ----------------------------------------------------------
// JWT HELPERS
// ----------------------------------------------------------

func (s *authService) createToken(userID uint, secret string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *authService) validateToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// ----------------------------------------------------------GetAuthUser
// GET AUTHENTICATED USER
func (s *authService) GetAuthUser(ctx context.Context) (*models.User, error) {
	// Get the user ID from the context
	userID, ok := ctx.Value("AuthID").(uint)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}

	// Find the user by ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found in the database")
	}

	return user, nil // Return a pointer to the user
}

// ----------------------------------------------------------Tracking User Online/Offline Status
// ----------------------------------------------------------

func (s *authService) TrackUserLogin(userID uint) error {
	// Create a unique key for the user session
	userKey := fmt.Sprintf("user:%d:online", userID)

	// Set the key in Redis to mark the user as online
	// The second parameter is the expiration time (e.g., 15 minutes)
	err := s.rdb.Set(context.Background(), userKey, "online", 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("could not track user login: %v", err)
	}

	return nil
}

// Track user logout by removing the Redis key
func (s *authService) TrackUserLogout(userID uint) error {
	// Create the same user key as used in login
	userKey := fmt.Sprintf("user:%d:online", userID)

	// Remove the key from Redis to mark the user as offline
	err := s.rdb.Del(context.Background(), userKey).Err()
	if err != nil {
		return fmt.Errorf("could not track user logout: %v", err)
	}

	return nil
}

// Check if a user is online
func (s *authService) IsUserOnline(userID uint) (bool, error) {
	// Create the user key to check
	userKey := fmt.Sprintf("user:%d:online", userID)

	// Check if the key exists in Redis
	val, err := s.rdb.Get(context.Background(), userKey).Result()
	if err == redis.Nil {
		// If the key doesn't exist, the user is offline
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check if user is online: %v", err)
	}

	// If the value is found, user is online
	return val == "online", nil
}
