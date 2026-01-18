# Head-to-Head Challenge Systems: Best Practices Research

## 1. Challenge Invitation/Acceptance Flows

### Standard Flow Patterns

**Direct Challenge Flow:**
- Challenger selects opponent → Creates challenge → Opponent receives notification → Accept/Decline → Match begins
- Used by: Chess.com, ESPN Fantasy, DraftKings

**Open Challenge Flow:**
- Challenger creates open challenge → Posted to challenge board → First to accept gets matched
- Used by: Lichess, Yahoo Fantasy

**Bracket/Tournament Integration:**
- System auto-generates challenges based on bracket progression
- Used by: March Madness brackets, FIFA tournaments

### Key Implementation Considerations

**Challenge Expiration:**
- Industry standard: 24-48 hours for direct challenges
- Open challenges: 1-7 days depending on platform activity
- Auto-decline expired challenges to prevent stale states

**Challenge Constraints:**
- Skill-based matching (ELO ranges, tier restrictions)
- Stake/entry fee matching for monetary competitions
- Friend-only vs public challenge options
- Rate limiting: Max 5-10 pending outgoing challenges per user

**State Management:**
```
PENDING → ACCEPTED → ACTIVE → COMPLETED
       ↘ DECLINED
       ↘ EXPIRED
       ↘ WITHDRAWN (before acceptance)
```

## 2. Scoring Mechanisms for 1v1 Competitions

### Sports Prediction Scoring Models

**Point-Based Systems:**
- Exact score: 5 points
- Correct winner + margin: 3 points  
- Correct winner only: 1 point
- Bonus multipliers for upset predictions

**Confidence-Based Scoring:**
- Users allocate confidence points (1-10) to predictions
- Higher confidence = higher reward/penalty
- Prevents conservative play, encourages strategic risk

**Progressive Scoring:**
- Early predictions worth more points
- Difficulty multipliers (playoff games worth 2x)
- Streak bonuses for consecutive correct predictions

### ELO Rating Integration
- Winner gains rating points based on opponent's rating
- Rating difference determines point exchange (10-50 points typical)
- Prevents farming weak opponents for easy wins

## 3. Notification Patterns for Challenge Updates

### Critical Notification Events

**High Priority (Push + In-App + Email):**
- Challenge received
- Challenge accepted/declined
- Match result determined
- Payment/reward distributed

**Medium Priority (Push + In-App):**
- Match starting soon (15 min warning)
- Opponent made prediction
- Leaderboard position change

**Low Priority (In-App only):**
- Challenge expired
- New open challenges available
- Weekly/monthly stats summary

### Notification Timing Strategy

**Real-time Events:**
- Challenge acceptance/decline: Immediate
- Live score updates: Every major event
- Match completion: Immediate

**Batched Events:**
- Daily digest of pending challenges
- Weekly performance summaries
- Monthly leaderboard updates

### Platform-Specific Considerations

**Mobile Push:**
- Respect quiet hours (10 PM - 8 AM local time)
- Frequency caps: Max 5 notifications per day
- Rich notifications with action buttons (Accept/Decline)

**Email:**
- HTML templates with clear CTAs
- Unsubscribe options for each notification type
- Mobile-responsive design

## 4. Database Schema Design for Challenges

### Core Tables

**challenges table:**
```sql
CREATE TABLE challenges (
    id UUID PRIMARY KEY,
    challenger_id UUID NOT NULL,
    opponent_id UUID,
    contest_id UUID NOT NULL,
    challenge_type ENUM('direct', 'open', 'bracket'),
    status ENUM('pending', 'accepted', 'active', 'completed', 'declined', 'expired', 'withdrawn'),
    stake_amount DECIMAL(10,2),
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
    accepted_at TIMESTAMP,
    completed_at TIMESTAMP,
    winner_id UUID,
    challenger_score INT DEFAULT 0,
    opponent_score INT DEFAULT 0,
    metadata JSONB
);
```

**challenge_predictions table:**
```sql
CREATE TABLE challenge_predictions (
    id UUID PRIMARY KEY,
    challenge_id UUID NOT NULL,
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    prediction JSONB NOT NULL,
    confidence_level INT,
    points_earned INT DEFAULT 0,
    created_at TIMESTAMP,
    UNIQUE(challenge_id, user_id, event_id)
);
```

### Indexing Strategy
```sql
-- Performance-critical queries
CREATE INDEX idx_challenges_opponent_status ON challenges(opponent_id, status);
CREATE INDEX idx_challenges_created_expires ON challenges(created_at, expires_at);
CREATE INDEX idx_challenge_predictions_lookup ON challenge_predictions(challenge_id, user_id);
```

### Partitioning Considerations
- Partition challenges by created_at (monthly partitions)
- Archive completed challenges older than 1 year
- Keep active/pending challenges in hot storage

## 5. Common Edge Cases and Solutions

### Timeout Handling

**Challenge Acceptance Timeout:**
- Auto-decline after expiration period
- Notify challenger of timeout
- Return any staked funds immediately
- Log timeout reason for analytics

**Prediction Deadline Timeout:**
- Lock in default predictions (if configured)
- Or forfeit the specific event prediction
- Partial scoring for submitted predictions only

**Implementation:**
```go
// Minimal timeout handler
func HandleChallengeTimeout(challengeID string) {
    challenge := getChallengeByID(challengeID)
    if challenge.Status == "pending" {
        challenge.Status = "expired"
        challenge.UpdatedAt = time.Now()
        updateChallenge(challenge)
        notifyChallenger(challenge.ChallengerID, "challenge_expired")
    }
}
```

### Tie Scenarios

**Exact Score Ties:**
- Declare both users winners (split rewards)
- Use tiebreaker criteria (prediction timestamps, confidence levels)
- Sudden death round with bonus prediction

**ELO Rating Ties:**
- Minimal rating exchange (±1-2 points)
- No rating change for both players
- Award participation points only

### Withdrawal Handling

**Before Acceptance:**
- Allow free withdrawal
- Notify potential opponents
- Remove from open challenge boards

**After Acceptance:**
- Forfeit penalty (lose staked amount)
- Opponent wins by default
- ELO rating penalty for withdrawer
- Cooldown period before creating new challenges

**Implementation:**
```go
func WithdrawChallenge(challengeID, userID string) error {
    challenge := getChallengeByID(challengeID)
    
    if challenge.Status == "pending" && challenge.ChallengerID == userID {
        challenge.Status = "withdrawn"
        return updateChallenge(challenge)
    }
    
    if challenge.Status == "accepted" {
        // Apply penalties
        applyWithdrawalPenalty(userID, challenge.StakeAmount)
        declareWinner(challenge, getOpponentID(challenge, userID))
        return updateChallenge(challenge)
    }
    
    return errors.New("cannot withdraw challenge in current state")
}
```

### Data Consistency Issues

**Concurrent Acceptance:**
- Use database transactions with row locking
- First successful transaction wins
- Notify other users of unavailable challenge

**Score Calculation Errors:**
- Implement idempotent scoring functions
- Store calculation audit trail
- Manual review process for disputed results

**Network Failures:**
- Implement retry mechanisms with exponential backoff
- Store interim states for recovery
- Graceful degradation for non-critical features

## Implementation Recommendations

### Phase 1: MVP Features
1. Direct challenge flow only
2. Simple point-based scoring
3. Basic push notifications
4. Core database schema

### Phase 2: Enhanced Features
1. Open challenge boards
2. ELO rating system
3. Rich notification templates
4. Advanced scoring mechanisms

### Phase 3: Advanced Features
1. Tournament bracket integration
2. Confidence-based scoring
3. Real-time updates via WebSocket
4. Comprehensive analytics dashboard

### Technology Stack Recommendations

**Backend:**
- Go with gRPC for service communication
- PostgreSQL for transactional data
- Redis for caching and real-time features
- Message queue (RabbitMQ/Kafka) for notifications

**Real-time Updates:**
- WebSocket connections for live updates
- Server-sent events for simpler implementations
- Push notification services (FCM, APNS)

**Monitoring:**
- Track challenge completion rates
- Monitor notification delivery success
- Alert on unusual timeout patterns
- Performance metrics for database queries

This research provides a comprehensive foundation for implementing robust head-to-head challenge systems that can scale and handle the complexities of competitive gaming and sports prediction platforms.