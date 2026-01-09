package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache provides Redis caching operations for leaderboards
type RedisCache struct {
	client *redis.Client
}

// Config holds Redis configuration
type Config struct {
	Addr            string
	Password        string
	DB              int
	ConnectTimeout  time.Duration
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(config Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection with configurable timeout
	timeout := config.ConnectTimeout
	if timeout == 0 {
		timeout = 5 * time.Second // Default timeout
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: rdb}, nil
}

// NewRedisCacheFromURL creates a new Redis cache client from URL
func NewRedisCacheFromURL(url string) (*RedisCache, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	// Test connection with default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: rdb}, nil
}

// LeaderboardEntry represents a leaderboard entry in cache
type LeaderboardEntry struct {
	UserID      uint    `json:"user_id"`
	TotalPoints float64 `json:"total_points"`
	Rank        uint    `json:"rank"`
}

// SetLeaderboardScore adds or updates a user's score in the contest leaderboard
func (r *RedisCache) SetLeaderboardScore(ctx context.Context, contestID uint, userID uint, points float64) error {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	member := strconv.FormatUint(uint64(userID), 10)
	
	return r.client.ZAdd(ctx, key, &redis.Z{
		Score:  points,
		Member: member,
	}).Err()
}

// GetLeaderboard retrieves the top N users from a contest leaderboard
func (r *RedisCache) GetLeaderboard(ctx context.Context, contestID uint, limit int64) ([]LeaderboardEntry, error) {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	
	// Get top scores in descending order
	results, err := r.client.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	entries := make([]LeaderboardEntry, len(results))
	for i, result := range results {
		member, ok := result.Member.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert member to string: %v", result.Member)
		}
		
		userID, err := strconv.ParseUint(member, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}

		entries[i] = LeaderboardEntry{
			UserID:      uint(userID),
			TotalPoints: result.Score,
			Rank:        uint(i + 1),
		}
	}

	return entries, nil
}

// GetUserRank retrieves a user's rank in a contest leaderboard
func (r *RedisCache) GetUserRank(ctx context.Context, contestID uint, userID uint) (uint, error) {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	member := strconv.FormatUint(uint64(userID), 10)
	
	rank, err := r.client.ZRevRank(ctx, key, member).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // User not found in leaderboard
		}
		return 0, fmt.Errorf("failed to get user rank: %w", err)
	}

	return uint(rank + 1), nil // Redis ranks are 0-based, we want 1-based
}

// GetUserScore retrieves a user's score in a contest leaderboard
func (r *RedisCache) GetUserScore(ctx context.Context, contestID uint, userID uint) (float64, error) {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	member := strconv.FormatUint(uint64(userID), 10)
	
	score, err := r.client.ZScore(ctx, key, member).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // User not found in leaderboard
		}
		return 0, fmt.Errorf("failed to get user score: %w", err)
	}

	return score, nil
}

// RemoveUserFromLeaderboard removes a user from a contest leaderboard
func (r *RedisCache) RemoveUserFromLeaderboard(ctx context.Context, contestID uint, userID uint) error {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	member := strconv.FormatUint(uint64(userID), 10)
	
	return r.client.ZRem(ctx, key, member).Err()
}

// GetLeaderboardSize returns the number of users in a contest leaderboard
func (r *RedisCache) GetLeaderboardSize(ctx context.Context, contestID uint) (int64, error) {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	return r.client.ZCard(ctx, key).Result()
}

// ClearLeaderboard removes all entries from a contest leaderboard
func (r *RedisCache) ClearLeaderboard(ctx context.Context, contestID uint) error {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	return r.client.Del(ctx, key).Err()
}

// SetLeaderboardTTL sets expiration time for a contest leaderboard
func (r *RedisCache) SetLeaderboardTTL(ctx context.Context, contestID uint, ttl time.Duration) error {
	key := fmt.Sprintf("contest:%d:leaderboard", contestID)
	return r.client.Expire(ctx, key, ttl).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
