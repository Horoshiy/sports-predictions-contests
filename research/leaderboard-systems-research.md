# Leaderboard Systems Research: Best Practices and Implementation Patterns

## Executive Summary

This research document covers comprehensive leaderboard system design for sports prediction platforms, focusing on real-time updates, ranking algorithms, caching strategies, and performance optimization techniques.

## 1. Real-Time Leaderboard Architecture Patterns

### 1.1 Event-Driven Architecture

```go
// Event-driven leaderboard updates
type LeaderboardEvent struct {
    Type      string    `json:"type"`      // "prediction_scored", "contest_updated"
    UserID    uint      `json:"user_id"`
    ContestID uint      `json:"contest_id"`
    Points    int       `json:"points"`
    Timestamp time.Time `json:"timestamp"`
}

type LeaderboardEventHandler struct {
    redis    *redis.Client
    db       *gorm.DB
    pubsub   *redis.PubSub
}

func (h *LeaderboardEventHandler) HandlePredictionScored(event LeaderboardEvent) error {
    // Update user's total score in Redis
    key := fmt.Sprintf("contest:%d:user:%d:score", event.ContestID, event.UserID)
    newScore, err := h.redis.IncrBy(key, int64(event.Points)).Result()
    if err != nil {
        return err
    }
    
    // Update sorted set for leaderboard ranking
    leaderboardKey := fmt.Sprintf("leaderboard:contest:%d", event.ContestID)
    return h.redis.ZAdd(leaderboardKey, &redis.Z{
        Score:  float64(newScore),
        Member: event.UserID,
    }).Err()
}
```

### 1.2 WebSocket Real-Time Updates

```go
type LeaderboardWebSocket struct {
    hub        *WebSocketHub
    redis      *redis.Client
    subscriber *redis.PubSub
}

type WebSocketHub struct {
    clients    map[*WebSocketClient]bool
    broadcast  chan []byte
    register   chan *WebSocketClient
    unregister chan *WebSocketClient
}

func (ws *LeaderboardWebSocket) BroadcastLeaderboardUpdate(contestID uint, update LeaderboardUpdate) {
    message := LeaderboardMessage{
        Type:      "leaderboard_update",
        ContestID: contestID,
        Data:      update,
    }
    
    ws.hub.broadcast <- message.ToJSON()
}

type LeaderboardUpdate struct {
    UserID      uint   `json:"user_id"`
    Username    string `json:"username"`
    NewScore    int    `json:"new_score"`
    OldRank     int    `json:"old_rank"`
    NewRank     int    `json:"new_rank"`
    PointsGain  int    `json:"points_gain"`
}
```

## 2. Ranking Algorithms and Strategies

### 2.1 Standard Ranking Systems

```go
type RankingStrategy interface {
    CalculateRanks(scores []UserScore) []RankedUser
    HandleTies(users []RankedUser) []RankedUser
}

// Standard competition ranking (1224 ranking)
type StandardRanking struct{}

func (sr *StandardRanking) CalculateRanks(scores []UserScore) []RankedUser {
    // Sort by score descending, then by tiebreaker criteria
    sort.Slice(scores, func(i, j int) bool {
        if scores[i].TotalPoints == scores[j].TotalPoints {
            return sr.resolveTie(scores[i], scores[j])
        }
        return scores[i].TotalPoints > scores[j].TotalPoints
    })
    
    ranks := make([]RankedUser, len(scores))
    currentRank := 1
    
    for i, score := range scores {
        if i > 0 && scores[i-1].TotalPoints != score.TotalPoints {
            currentRank = i + 1
        }
        
        ranks[i] = RankedUser{
            UserScore: score,
            Rank:      currentRank,
        }
    }
    
    return ranks
}

func (sr *StandardRanking) resolveTie(a, b UserScore) bool {
    // Tiebreaker 1: Higher accuracy
    if a.Accuracy != b.Accuracy {
        return a.Accuracy > b.Accuracy
    }
    
    // Tiebreaker 2: Earlier submission time for latest prediction
    return a.LastSubmissionTime.Before(b.LastSubmissionTime)
}
```

### 2.2 ELO-Based Dynamic Ranking

```go
type ELORanking struct {
    KFactor    float64 // Rating change sensitivity
    BaseRating float64 // Starting rating
}

func (elo *ELORanking) UpdateRating(userRating, opponentRating float64, result float64) float64 {
    expectedScore := 1.0 / (1.0 + math.Pow(10, (opponentRating-userRating)/400))
    return userRating + elo.KFactor*(result-expectedScore)
}

func (elo *ELORanking) CalculateContestRatings(predictions []PredictionResult) map[uint]float64 {
    ratings := make(map[uint]float64)
    
    // Initialize ratings
    for _, pred := range predictions {
        if _, exists := ratings[pred.UserID]; !exists {
            ratings[pred.UserID] = elo.BaseRating
        }
    }
    
    // Process each event's predictions
    eventGroups := groupPredictionsByEvent(predictions)
    
    for _, eventPreds := range eventGroups {
        avgRating := calculateAverageRating(eventPreds, ratings)
        
        for _, pred := range eventPreds {
            result := 0.0
            if pred.IsCorrect {
                result = 1.0
            }
            
            ratings[pred.UserID] = elo.UpdateRating(
                ratings[pred.UserID],
                avgRating,
                result,
            )
        }
    }
    
    return ratings
}
```

### 2.3 Percentile-Based Ranking

```go
type PercentileRanking struct {
    Buckets int // Number of percentile buckets
}

func (pr *PercentileRanking) CalculatePercentiles(scores []UserScore) []PercentileRank {
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].TotalPoints < scores[j].TotalPoints
    })
    
    ranks := make([]PercentileRank, len(scores))
    
    for i, score := range scores {
        percentile := float64(i) / float64(len(scores)) * 100
        bucket := int(percentile / (100.0 / float64(pr.Buckets)))
        
        ranks[i] = PercentileRank{
            UserScore:  score,
            Percentile: percentile,
            Bucket:     bucket,
        }
    }
    
    return ranks
}
```

## 3. Caching Strategies for High Performance

### 3.1 Multi-Layer Caching Architecture

```go
type LeaderboardCache struct {
    L1Cache *sync.Map          // In-memory cache
    L2Cache *redis.Client      // Redis cache
    L3Cache *gorm.DB          // Database
    
    TTL     time.Duration
    MaxSize int
}

func (lc *LeaderboardCache) GetLeaderboard(contestID uint, limit int) ([]RankedUser, error) {
    // L1: Check in-memory cache
    if cached, ok := lc.L1Cache.Load(fmt.Sprintf("leaderboard:%d", contestID)); ok {
        if leaderboard, ok := cached.(CachedLeaderboard); ok && !leaderboard.IsExpired() {
            return leaderboard.Data[:min(limit, len(leaderboard.Data))], nil
        }
    }
    
    // L2: Check Redis cache
    leaderboard, err := lc.getFromRedis(contestID, limit)
    if err == nil {
        lc.setInMemory(contestID, leaderboard)
        return leaderboard, nil
    }
    
    // L3: Fallback to database
    leaderboard, err = lc.getFromDatabase(contestID, limit)
    if err != nil {
        return nil, err
    }
    
    // Cache in both layers
    go lc.setInRedis(contestID, leaderboard)
    lc.setInMemory(contestID, leaderboard)
    
    return leaderboard, nil
}
```

### 3.2 Redis Sorted Sets for Leaderboards

```go
type RedisLeaderboard struct {
    client *redis.Client
}

func (rl *RedisLeaderboard) UpdateUserScore(contestID, userID uint, points int) error {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    
    // Use pipeline for atomic operations
    pipe := rl.client.Pipeline()
    
    // Update score in sorted set
    pipe.ZIncrBy(key, float64(points), fmt.Sprintf("user:%d", userID))
    
    // Update user metadata
    userKey := fmt.Sprintf("user:%d:contest:%d", userID, contestID)
    pipe.HIncrBy(userKey, "total_points", int64(points))
    pipe.HIncrBy(userKey, "predictions_count", 1)
    
    // Set expiration
    pipe.Expire(key, 24*time.Hour)
    pipe.Expire(userKey, 24*time.Hour)
    
    _, err := pipe.Exec()
    return err
}

func (rl *RedisLeaderboard) GetTopUsers(contestID uint, limit int) ([]RankedUser, error) {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    
    // Get top users with scores
    results, err := rl.client.ZRevRangeWithScores(key, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, err
    }
    
    users := make([]RankedUser, len(results))
    for i, result := range results {
        userID := extractUserID(result.Member.(string))
        users[i] = RankedUser{
            UserID: userID,
            Score:  int(result.Score),
            Rank:   i + 1,
        }
    }
    
    return users, nil
}

func (rl *RedisLeaderboard) GetUserRank(contestID, userID uint) (int, error) {
    key := fmt.Sprintf("leaderboard:contest:%d", contestID)
    member := fmt.Sprintf("user:%d", userID)
    
    rank, err := rl.client.ZRevRank(key, member).Result()
    if err != nil {
        return 0, err
    }
    
    return int(rank) + 1, nil // Redis ranks are 0-based
}
```

### 3.3 Cache Invalidation Strategies

```go
type CacheInvalidator struct {
    redis     *redis.Client
    memCache  *sync.Map
    pubsub    *redis.PubSub
}

func (ci *CacheInvalidator) InvalidateLeaderboard(contestID uint) error {
    // Invalidate Redis cache
    pattern := fmt.Sprintf("leaderboard:contest:%d*", contestID)
    keys, err := ci.redis.Keys(pattern).Result()
    if err != nil {
        return err
    }
    
    if len(keys) > 0 {
        ci.redis.Del(keys...)
    }
    
    // Invalidate memory cache
    ci.memCache.Delete(fmt.Sprintf("leaderboard:%d", contestID))
    
    // Publish invalidation event
    event := CacheInvalidationEvent{
        Type:      "leaderboard_invalidated",
        ContestID: contestID,
        Timestamp: time.Now(),
    }
    
    return ci.redis.Publish("cache_invalidation", event.ToJSON()).Err()
}

func (ci *CacheInvalidator) HandleInvalidationEvents() {
    ch := ci.pubsub.Channel()
    
    for msg := range ch {
        var event CacheInvalidationEvent
        if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
            continue
        }
        
        switch event.Type {
        case "leaderboard_invalidated":
            ci.memCache.Delete(fmt.Sprintf("leaderboard:%d", event.ContestID))
        }
    }
}
```

## 4. Performance Optimization Techniques

### 4.1 Batch Processing for Score Updates

```go
type BatchScoreProcessor struct {
    batchSize    int
    flushInterval time.Duration
    buffer       []ScoreUpdate
    mutex        sync.Mutex
    redis        *redis.Client
}

type ScoreUpdate struct {
    ContestID uint
    UserID    uint
    Points    int
    Timestamp time.Time
}

func (bsp *BatchScoreProcessor) AddScoreUpdate(update ScoreUpdate) {
    bsp.mutex.Lock()
    defer bsp.mutex.Unlock()
    
    bsp.buffer = append(bsp.buffer, update)
    
    if len(bsp.buffer) >= bsp.batchSize {
        go bsp.flushBuffer()
    }
}

func (bsp *BatchScoreProcessor) flushBuffer() {
    bsp.mutex.Lock()
    updates := make([]ScoreUpdate, len(bsp.buffer))
    copy(updates, bsp.buffer)
    bsp.buffer = bsp.buffer[:0]
    bsp.mutex.Unlock()
    
    // Group updates by contest
    contestUpdates := make(map[uint][]ScoreUpdate)
    for _, update := range updates {
        contestUpdates[update.ContestID] = append(contestUpdates[update.ContestID], update)
    }
    
    // Process each contest's updates in a pipeline
    for contestID, updates := range contestUpdates {
        bsp.processBatch(contestID, updates)
    }
}

func (bsp *BatchScoreProcessor) processBatch(contestID uint, updates []ScoreUpdate) error {
    pipe := bsp.redis.Pipeline()
    leaderboardKey := fmt.Sprintf("leaderboard:contest:%d", contestID)
    
    for _, update := range updates {
        member := fmt.Sprintf("user:%d", update.UserID)
        pipe.ZIncrBy(leaderboardKey, float64(update.Points), member)
    }
    
    _, err := pipe.Exec()
    return err
}
```

### 4.2 Precomputed Leaderboard Snapshots

```go
type LeaderboardSnapshot struct {
    ContestID   uint      `json:"contest_id"`
    Data        []byte    `json:"data"`        // Compressed leaderboard data
    Version     int       `json:"version"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   time.Time `json:"expires_at"`
}

type SnapshotManager struct {
    db    *gorm.DB
    redis *redis.Client
}

func (sm *SnapshotManager) CreateSnapshot(contestID uint) error {
    // Get current leaderboard
    leaderboard, err := sm.getCurrentLeaderboard(contestID)
    if err != nil {
        return err
    }
    
    // Compress data
    compressed, err := sm.compressLeaderboard(leaderboard)
    if err != nil {
        return err
    }
    
    snapshot := LeaderboardSnapshot{
        ContestID: contestID,
        Data:      compressed,
        Version:   sm.getNextVersion(contestID),
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(1 * time.Hour),
    }
    
    // Store in database
    if err := sm.db.Create(&snapshot).Error; err != nil {
        return err
    }
    
    // Cache in Redis
    key := fmt.Sprintf("snapshot:contest:%d:latest", contestID)
    return sm.redis.Set(key, snapshot.Data, time.Hour).Err()
}

func (sm *SnapshotManager) GetSnapshot(contestID uint) ([]RankedUser, error) {
    key := fmt.Sprintf("snapshot:contest:%d:latest", contestID)
    
    data, err := sm.redis.Get(key).Bytes()
    if err != nil {
        // Fallback to database
        var snapshot LeaderboardSnapshot
        if err := sm.db.Where("contest_id = ? AND expires_at > ?", contestID, time.Now()).
            Order("version DESC").First(&snapshot).Error; err != nil {
            return nil, err
        }
        data = snapshot.Data
    }
    
    return sm.decompressLeaderboard(data)
}
```

### 4.3 Database Query Optimization

```sql
-- Optimized leaderboard query with proper indexing
CREATE INDEX CONCURRENTLY idx_user_contest_scores_leaderboard 
ON user_contest_scores (contest_id, total_points DESC, accuracy_percentage DESC, last_updated);

-- Materialized view for complex leaderboard calculations
CREATE MATERIALIZED VIEW contest_leaderboard_mv AS
SELECT 
    ucs.contest_id,
    ucs.user_id,
    u.username,
    u.display_name,
    ucs.total_points,
    ucs.correct_predictions,
    ucs.total_predictions,
    ucs.accuracy_percentage,
    ROW_NUMBER() OVER (
        PARTITION BY ucs.contest_id 
        ORDER BY ucs.total_points DESC, 
                 ucs.accuracy_percentage DESC, 
                 ucs.last_updated ASC
    ) as rank
FROM user_contest_scores ucs
JOIN users u ON u.id = ucs.user_id
WHERE ucs.total_predictions > 0;

-- Refresh strategy
CREATE OR REPLACE FUNCTION refresh_leaderboard_mv()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY contest_leaderboard_mv;
END;
$$ LANGUAGE plpgsql;
```

## 5. Scalability Patterns

### 5.1 Horizontal Partitioning

```go
type PartitionedLeaderboard struct {
    partitions map[string]*LeaderboardPartition
    hasher     consistent.Consistent
}

type LeaderboardPartition struct {
    ID     string
    Redis  *redis.Client
    Weight int
}

func (pl *PartitionedLeaderboard) getPartition(contestID uint) *LeaderboardPartition {
    key := fmt.Sprintf("contest:%d", contestID)
    partitionID := pl.hasher.Get(key)
    return pl.partitions[partitionID]
}

func (pl *PartitionedLeaderboard) UpdateScore(contestID, userID uint, points int) error {
    partition := pl.getPartition(contestID)
    return partition.UpdateScore(contestID, userID, points)
}

func (pl *PartitionedLeaderboard) GetLeaderboard(contestID uint, limit int) ([]RankedUser, error) {
    partition := pl.getPartition(contestID)
    return partition.GetTopUsers(contestID, limit)
}
```

### 5.2 Read Replicas for Leaderboard Queries

```go
type LeaderboardService struct {
    writeDB *gorm.DB
    readDB  *gorm.DB
    redis   *redis.Client
}

func (ls *LeaderboardService) GetLeaderboard(contestID uint, limit int) ([]RankedUser, error) {
    // Try cache first
    if cached, err := ls.getFromCache(contestID, limit); err == nil {
        return cached, nil
    }
    
    // Use read replica for database queries
    var scores []UserContestScore
    err := ls.readDB.Where("contest_id = ?", contestID).
        Order("total_points DESC, accuracy_percentage DESC").
        Limit(limit).
        Find(&scores).Error
    
    if err != nil {
        return nil, err
    }
    
    return ls.convertToRankedUsers(scores), nil
}

func (ls *LeaderboardService) UpdateScore(contestID, userID uint, points int) error {
    // Use write database for updates
    return ls.writeDB.Model(&UserContestScore{}).
        Where("contest_id = ? AND user_id = ?", contestID, userID).
        Update("total_points", gorm.Expr("total_points + ?", points)).Error
}
```

## 6. Real-Time Features Implementation

### 6.1 Live Score Updates

```go
type LiveScoreUpdater struct {
    websocket *WebSocketManager
    redis     *redis.Client
    eventBus  *EventBus
}

func (lsu *LiveScoreUpdater) HandleEventResult(eventID uint, result EventResult) error {
    // Get all predictions for this event
    predictions, err := lsu.getPredictionsForEvent(eventID)
    if err != nil {
        return err
    }
    
    // Calculate scores and update leaderboards
    for _, prediction := range predictions {
        score := lsu.calculateScore(prediction, result)
        
        // Update user's total score
        if err := lsu.updateUserScore(prediction.ContestID, prediction.UserID, score); err != nil {
            log.Printf("Failed to update score: %v", err)
            continue
        }
        
        // Broadcast real-time update
        update := LiveScoreUpdate{
            UserID:      prediction.UserID,
            ContestID:   prediction.ContestID,
            EventID:     eventID,
            PointsGain:  score,
            NewRank:     lsu.getUserRank(prediction.ContestID, prediction.UserID),
        }
        
        lsu.websocket.BroadcastToContest(prediction.ContestID, update)
    }
    
    return nil
}
```

### 6.2 Leaderboard Animations and Transitions

```typescript
// Frontend leaderboard animation system
class LeaderboardAnimator {
    private container: HTMLElement;
    private users: Map<number, LeaderboardUser> = new Map();
    
    updateLeaderboard(newData: LeaderboardUser[]) {
        const oldPositions = this.getCurrentPositions();
        
        // Calculate position changes
        const changes = this.calculateChanges(oldPositions, newData);
        
        // Animate transitions
        this.animateChanges(changes);
    }
    
    private animateChanges(changes: PositionChange[]) {
        changes.forEach(change => {
            const element = this.getUserElement(change.userId);
            
            if (change.type === 'rank_up') {
                this.animateRankUp(element, change.oldRank, change.newRank);
            } else if (change.type === 'rank_down') {
                this.animateRankDown(element, change.oldRank, change.newRank);
            } else if (change.type === 'score_update') {
                this.animateScoreUpdate(element, change.pointsGain);
            }
        });
    }
    
    private animateRankUp(element: HTMLElement, oldRank: number, newRank: number) {
        element.classList.add('rank-up');
        
        // Smooth transition animation
        element.style.transform = `translateY(${(oldRank - newRank) * 60}px)`;
        
        setTimeout(() => {
            element.style.transform = 'translateY(0)';
            element.classList.remove('rank-up');
        }, 300);
    }
}
```

## 7. Implementation Recommendations

### 7.1 Technology Stack

**Backend:**
- **Database**: PostgreSQL with read replicas
- **Cache**: Redis Cluster for distributed caching
- **Message Queue**: Redis Pub/Sub or Apache Kafka
- **WebSockets**: Socket.io or native WebSocket with connection pooling

**Frontend:**
- **Real-time Updates**: WebSocket connections with reconnection logic
- **State Management**: Redux/Zustand for leaderboard state
- **Animations**: Framer Motion or CSS transitions

### 7.2 Performance Benchmarks

Target performance metrics:
- **Leaderboard Query Response**: < 100ms for top 100 users
- **Real-time Update Latency**: < 500ms from score calculation to UI update
- **Cache Hit Rate**: > 95% for leaderboard queries
- **Concurrent Users**: Support 10,000+ concurrent leaderboard viewers

### 7.3 Monitoring and Alerting

```go
type LeaderboardMetrics struct {
    QueryLatency     prometheus.Histogram
    CacheHitRate     prometheus.Counter
    UpdateThroughput prometheus.Counter
    ErrorRate        prometheus.Counter
}

func (lm *LeaderboardMetrics) RecordQuery(duration time.Duration, cacheHit bool) {
    lm.QueryLatency.Observe(duration.Seconds())
    
    if cacheHit {
        lm.CacheHitRate.With(prometheus.Labels{"type": "hit"}).Inc()
    } else {
        lm.CacheHitRate.With(prometheus.Labels{"type": "miss"}).Inc()
    }
}
```

## 8. Security Considerations

### 8.1 Rate Limiting

```go
type LeaderboardRateLimiter struct {
    redis *redis.Client
}

func (lrl *LeaderboardRateLimiter) CheckRateLimit(userID uint, action string) error {
    key := fmt.Sprintf("rate_limit:%s:%d", action, userID)
    
    current, err := lrl.redis.Incr(key).Result()
    if err != nil {
        return err
    }
    
    if current == 1 {
        lrl.redis.Expire(key, time.Minute)
    }
    
    limit := lrl.getLimit(action)
    if current > limit {
        return errors.New("rate limit exceeded")
    }
    
    return nil
}
```

### 8.2 Data Validation and Sanitization

```go
func (ls *LeaderboardService) ValidateLeaderboardRequest(req LeaderboardRequest) error {
    if req.ContestID == 0 {
        return errors.New("invalid contest ID")
    }
    
    if req.Limit < 1 || req.Limit > 1000 {
        return errors.New("limit must be between 1 and 1000")
    }
    
    if req.Offset < 0 {
        return errors.New("offset cannot be negative")
    }
    
    return nil
}
```

This comprehensive research provides the foundation for implementing a high-performance, scalable leaderboard system for the sports prediction platform.