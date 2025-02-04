package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds the Redis connection details
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// SetupRedis initializes a Redis client
func SetupRedis(ctx context.Context, cfg *RedisConfig) (*redis.Client, error) {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping to test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}

// LoadRedisConfig from environment variables
func LoadRedisConfig() *RedisConfig {
	// Example environment variables: REDIS_HOST, REDIS_PORT, REDIS_PASS, REDIS_DB
	host := getEnv("REDIS_HOST", "127.0.0.1")
	port := getEnv("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASS") // empty if not set
	db := 0
	// parse DB from environment if needed (omitted here for brevity)

	return &RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}
