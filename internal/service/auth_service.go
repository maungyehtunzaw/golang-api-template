package service

import (
	"context"
	"fmt"
	"time"

	"golang-api-template/internal/config"
	"golang-api-template/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(refreshToken string) error
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
func (s *authService) Login(email, password string) (string, string, error) {
	// 1. Find user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// 2. Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// 3. Create access token
	accessToken, err := s.createToken(user.ID, s.cfg.JWTAccessSecret, time.Minute*time.Duration(s.cfg.AccessTokenExpireMin))
	if err != nil {
		return "", "", err
	}

	// 4. Create refresh token
	refreshToken, err := s.createToken(user.ID, s.cfg.JWTRefreshSecret, time.Hour*time.Duration(s.cfg.RefreshTokenExpireHrs))
	if err != nil {
		return "", "", err
	}

	// 5. Store the refresh token in Redis
	ctx := context.Background()
	// We store the token as key -> userID, so we can reference it later
	err = s.rdb.Set(ctx, refreshToken, user.ID, time.Hour*time.Duration(s.cfg.RefreshTokenExpireHrs)).Err()
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
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
