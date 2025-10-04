package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"expensio-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// InitRedis initializes Redis connection
func InitRedis(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// CloseRedis closes Redis connection
func CloseRedis() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			fmt.Printf("Error closing Redis connection: %v\n", err)
		} else {
			fmt.Println("✅ Disconnected from Redis")
		}
	}
}

// Set stores a value in Redis with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return Client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from Redis
func Get(key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key from Redis
func Delete(keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// DeletePattern deletes all keys matching a pattern
func DeletePattern(pattern string) error {
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := Client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Exists checks if a key exists in Redis
func Exists(key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetString stores a string value in Redis
func SetString(key string, value string, ttl time.Duration) error {
	return Client.Set(ctx, key, value, ttl).Err()
}

// GetString retrieves a string value from Redis
func GetString(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Increment increments a value in Redis
func Increment(key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Expire sets TTL for a key
func Expire(key string, ttl time.Duration) error {
	return Client.Expire(ctx, key, ttl).Err()
}

// TTL gets remaining TTL for a key
func TTL(key string) (time.Duration, error) {
	return Client.TTL(ctx, key).Result()
}

// SetHash stores a hash in Redis
func SetHash(key string, field string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := Client.HSet(ctx, key, field, data).Err(); err != nil {
		return err
	}

	return Client.Expire(ctx, key, ttl).Err()
}

// GetHash retrieves a hash field from Redis
func GetHash(key string, field string, dest interface{}) error {
	data, err := Client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key or field not found")
		}
		return fmt.Errorf("failed to get hash value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// DeleteHashField removes a field from a hash
func DeleteHashField(key string, fields ...string) error {
	return Client.HDel(ctx, key, fields...).Err()
}

// GetAllHash retrieves all fields from a hash
func GetAllHash(key string) (map[string]string, error) {
	return Client.HGetAll(ctx, key).Result()
}

// AddToSet adds a member to a set
func AddToSet(key string, members ...interface{}) error {
	return Client.SAdd(ctx, key, members...).Err()
}

// RemoveFromSet removes a member from a set
func RemoveFromSet(key string, members ...interface{}) error {
	return Client.SRem(ctx, key, members...).Err()
}

// IsMemberOfSet checks if a member exists in a set
func IsMemberOfSet(key string, member interface{}) (bool, error) {
	return Client.SIsMember(ctx, key, member).Result()
}

// GetSetMembers retrieves all members from a set
func GetSetMembers(key string) ([]string, error) {
	return Client.SMembers(ctx, key).Result()
}

// HealthCheck checks if Redis connection is alive
func HealthCheck() error {
	return Client.Ping(ctx).Err()
}

// FlushDB clears all keys from current database (use with caution!)
func FlushDB() error {
	return Client.FlushDB(ctx).Err()
}
