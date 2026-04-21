package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host     string
	Port     string
	Password string
}

func Connect(cfg Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("Successfully connected to Redis")
	return client, nil
}

// GetAvailabilityKey returns a consistent Redis key pattern for equipment availability
func GetAvailabilityKey(tenantID, departmentID, categoryID string) string {
	return fmt.Sprintf("availability:%s:%s:%s", tenantID, departmentID, categoryID)
}

// GetAvailabilityBoardKey returns the key for the entire tenant availability board hash
func GetAvailabilityBoardKey(tenantID string) string {
	return fmt.Sprintf("availability_board:%s", tenantID)
}
