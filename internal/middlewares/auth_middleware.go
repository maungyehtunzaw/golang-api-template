package middlewares

import (
	"net/http"
	"strings"

	"golang-api-template/internal/config"
	"golang-api-template/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing Authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the Access token
		claims, err := validateAccessToken(tokenStr, cfg.JWTAccessSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired access token")
			c.Abort()
			return
		}

		// Optionally store user ID in context
		userID, ok := claims["user_id"].(float64)
		if ok {
			c.Set("userID", uint(userID))
		}

		c.Next()
	}
}

func validateAccessToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
