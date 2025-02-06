package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang-api-template/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	Port       string

	// JWT secrets & expirations
	JWTAccessSecret       string
	JWTRefreshSecret      string
	AccessTokenExpireMin  int
	RefreshTokenExpireHrs int
	RestTokenExpireInMin  int

	// ... possibly more fields
	Redis *RedisConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using default environment variables")
	}

	accessExp, _ := strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRE_MIN", "15"))
	refreshExp, _ := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRE_HOUR", "72"))
	resetTokenExpire, _ := strconv.Atoi(getEnv("RESET_TOKEN_EXPIRY_MIN", "15"))

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "demo"),
		DBPort:     getEnv("DB_PORT", "3306"),
		Port:       getEnv("PORT", "8080"),

		JWTAccessSecret:       getEnv("JWT_ACCESS_SECRET", "access-secret-example"),
		JWTRefreshSecret:      getEnv("JWT_REFRESH_SECRET", "refresh-secret-example"),
		RestTokenExpireInMin:  resetTokenExpire,
		AccessTokenExpireMin:  accessExp,
		RefreshTokenExpireHrs: refreshExp,

		Redis: LoadRedisConfig(), // from redis.go
	}
	return cfg, nil
}

func SetupDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
}

func GetResetTokenExpiry() time.Duration {
	// Convert the environment variable to an integer
	expiryMinutes, err := strconv.Atoi(getEnv("RESET_TOKEN_EXPIRY_MIN", "15"))
	if err != nil {
		// Handle error or fall back to a default value
		expiryMinutes = 15 // Default to 15 minutes if there's an error
	}
	return time.Minute * time.Duration(expiryMinutes)
}

// Helper
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
