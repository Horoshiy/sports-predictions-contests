# Minimal Leaderboard Implementation Example

This document provides minimal, focused code examples demonstrating the core leaderboard patterns researched.

## 1. Redis-Based Real-Time Leaderboard

```go
package leaderboard

import (
    "fmt"
    "strconv"
    "strings"
    "time"
    
    "github.com/go-redis/redis/v8"
    "context"
)

type RedisLeaderboard struct {
    client *redis.Client
}

func NewRedisLeaderboard(client *redis.Client) *RedisLeaderboard {
    return &RedisLeaderboard{client: client}
}

// UpdateScore atomically updates user score and rank
func (r *RedisLeaderboard) UpdateScore(ctx context.Context, contestID, userID uint, points int) error {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    member := fmt.Sprintf("user:%d", userID)
    
    return r.client.ZIncrBy(ctx, key, float64(points), member).Err()
}

// GetTopUsers returns top N users with their ranks
func (r *RedisLeaderboard) GetTopUsers(ctx context.Context, contestID uint, limit int) ([]RankedUser, error) {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    
    results, err := r.client.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, err
    }
    
    users := make([]RankedUser, len(results))
    for i, result := range results {
        userID, _ := strconv.ParseUint(strings.TrimPrefix(result.Member.(string), "user:"), 10, 32)
        users[i] = RankedUser{
            UserID: uint(userID),
            Score:  int(result.Score),
            Rank:   i + 1,
        }
    }
    
    return users, nil
}

// GetUserRank returns specific user's rank
func (r *RedisLeaderboard) GetUserRank(ctx context.Context, contestID, userID uint) (int, error) {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    member := fmt.Sprintf("user:%d", userID)
    
    rank, err := r.client.ZRevRank(ctx, key, member).Result()
    if err != nil {
        return 0, err
    }
    
    return int(rank) + 1, nil
}

type RankedUser struct {
    UserID uint `json:"user_id"`
    Score  int  `json:"score"`
    Rank   int  `json:"rank"`
}
```

## 2. Event-Driven Score Updates

```go
package events

import (
    "context"
    "encoding/json"
    "time"
)

type ScoreEvent struct {
    UserID    uint      `json:"user_id"`
    ContestID uint      `json:"contest_id"`
    Points    int       `json:"points"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
}

type EventProcessor struct {
    leaderboard *RedisLeaderboard
    pubsub      *redis.PubSub
}

func (ep *EventProcessor) ProcessScoreEvent(ctx context.Context, event ScoreEvent) error {
    // Update leaderboard
    if err := ep.leaderboard.UpdateScore(ctx, event.ContestID, event.UserID, event.Points); err != nil {
        return err
    }
    
    // Publish update for real-time clients
    update := LeaderboardUpdate{
        UserID:    event.UserID,
        ContestID: event.ContestID,
        Points:    event.Points,
        Timestamp: time.Now(),
    }
    
    data, _ := json.Marshal(update)
    return ep.pubsub.Publish(ctx, fmt.Sprintf("leaderboard:%d", event.ContestID), data).Err()
}

type LeaderboardUpdate struct {
    UserID    uint      `json:"user_id"`
    ContestID uint      `json:"contest_id"`
    Points    int       `json:"points"`
    Timestamp time.Time `json:"timestamp"`
}
```

## 3. Caching Layer

```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

type CachedLeaderboard struct {
    Data      []RankedUser `json:"data"`
    ExpiresAt time.Time    `json:"expires_at"`
}

func (c CachedLeaderboard) IsExpired() bool {
    return time.Now().After(c.ExpiresAt)
}

type LeaderboardCache struct {
    memory sync.Map
    redis  *redis.Client
    ttl    time.Duration
}

func NewLeaderboardCache(redis *redis.Client, ttl time.Duration) *LeaderboardCache {
    return &LeaderboardCache{
        redis: redis,
        ttl:   ttl,
    }
}

func (lc *LeaderboardCache) Get(ctx context.Context, contestID uint, limit int) ([]RankedUser, bool) {
    key := fmt.Sprintf("leaderboard:%d:%d", contestID, limit)
    
    // Check memory cache
    if cached, ok := lc.memory.Load(key); ok {
        if leaderboard, ok := cached.(CachedLeaderboard); ok && !leaderboard.IsExpired() {
            return leaderboard.Data, true
        }
    }
    
    // Check Redis cache
    data, err := lc.redis.Get(ctx, key).Result()
    if err == nil {
        var leaderboard CachedLeaderboard
        if json.Unmarshal([]byte(data), &leaderboard) == nil && !leaderboard.IsExpired() {
            lc.memory.Store(key, leaderboard)
            return leaderboard.Data, true
        }
    }
    
    return nil, false
}

func (lc *LeaderboardCache) Set(ctx context.Context, contestID uint, limit int, data []RankedUser) {
    key := fmt.Sprintf("leaderboard:%d:%d", contestID, limit)
    
    cached := CachedLeaderboard{
        Data:      data,
        ExpiresAt: time.Now().Add(lc.ttl),
    }
    
    // Store in memory
    lc.memory.Store(key, cached)
    
    // Store in Redis
    if jsonData, err := json.Marshal(cached); err == nil {
        lc.redis.Set(ctx, key, jsonData, lc.ttl)
    }
}
```

## 4. WebSocket Real-Time Updates

```go
package websocket

import (
    "encoding/json"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
)

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

type Client struct {
    hub       *Hub
    conn      *websocket.Conn
    send      chan []byte
    contestID uint
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mutex.Unlock()
            
        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

func (h *Hub) BroadcastToContest(contestID uint, update LeaderboardUpdate) {
    message, _ := json.Marshal(update)
    
    h.mutex.RLock()
    for client := range h.clients {
        if client.contestID == contestID {
            select {
            case client.send <- message:
            default:
                close(client.send)
                delete(h.clients, client)
            }
        }
    }
    h.mutex.RUnlock()
}
```

## 5. Complete Service Integration

```go
package service

import (
    "context"
    "time"
)

type LeaderboardService struct {
    redis     *RedisLeaderboard
    cache     *LeaderboardCache
    hub       *Hub
    processor *EventProcessor
}

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
    redis := NewRedisLeaderboard(redisClient)
    cache := NewLeaderboardCache(redisClient, 5*time.Minute)
    hub := NewHub()
    
    go hub.Run()
    
    return &LeaderboardService{
        redis: redis,
        cache: cache,
        hub:   hub,
        processor: &EventProcessor{
            leaderboard: redis,
        },
    }
}

func (ls *LeaderboardService) GetLeaderboard(ctx context.Context, contestID uint, limit int) ([]RankedUser, error) {
    // Try cache first
    if cached, found := ls.cache.Get(ctx, contestID, limit); found {
        return cached, nil
    }
    
    // Get from Redis
    users, err := ls.redis.GetTopUsers(ctx, contestID, limit)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    ls.cache.Set(ctx, contestID, limit, users)
    
    return users, nil
}

func (ls *LeaderboardService) UpdateScore(ctx context.Context, contestID, userID uint, points int) error {
    // Update score
    if err := ls.redis.UpdateScore(ctx, contestID, userID, points); err != nil {
        return err
    }
    
    // Invalidate cache
    ls.cache.Invalidate(contestID)
    
    // Broadcast update
    update := LeaderboardUpdate{
        UserID:    userID,
        ContestID: contestID,
        Points:    points,
        Timestamp: time.Now(),
    }
    
    ls.hub.BroadcastToContest(contestID, update)
    
    return nil
}
```

This minimal implementation demonstrates the core patterns for building a high-performance, real-time leaderboard system.