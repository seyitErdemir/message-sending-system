package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

type MessageCache struct {
	ID        uint   `json:"id"`
	MessageID string `json:"message_id"`
	Status    bool   `json:"status"`
	Content   string `json:"content"`
	Phone     string `json:"phone"`
}

func Connect() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	// Test connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

func SetMessageCache(messageID uint, data MessageCache) error {
	key := fmt.Sprintf("message:%d", messageID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %v", err)
	}

	// Store in cache for 1 hour
	err = RedisClient.Set(Ctx, key, jsonData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set message cache: %v", err)
	}

	return nil
}

func GetMessageCache(messageID uint) (*MessageCache, error) {
	key := fmt.Sprintf("message:%d", messageID)
	data, err := RedisClient.Get(Ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get message cache: %v", err)
	}

	var message MessageCache
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message data: %v", err)
	}

	return &message, nil
}

// GetMessageCacheWithTimeout gets a message from cache with a timeout
func GetMessageCacheWithTimeout(messageID uint) (*MessageCache, error) {
	// Create a new context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(Ctx, time.Second*5)
	// Clean up context when operation is done
	defer cancel()

	key := fmt.Sprintf("message:%d", messageID)

	data, err := RedisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err == context.DeadlineExceeded {
		return nil, fmt.Errorf("redis operation timed out after 5 seconds")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get message cache: %v", err)
	}

	var message MessageCache
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message data: %v", err)
	}

	return &message, nil
}
