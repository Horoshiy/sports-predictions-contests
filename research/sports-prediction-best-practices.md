# Sports Prediction System Best Practices & Data Models

## 1. Prediction Data Structures and Validation

### Core Prediction Structure
```go
type Prediction struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    ContestID   string    `json:"contest_id"`
    EventID     string    `json:"event_id"`
    Type        PredictionType `json:"type"`
    Value       interface{} `json:"value"`
    Confidence  float64   `json:"confidence,omitempty"`
    SubmittedAt time.Time `json:"submitted_at"`
    Status      PredictionStatus `json:"status"`
}

type PredictionType string
const (
    MatchOutcome    PredictionType = "match_outcome"
    ExactScore      PredictionType = "exact_score"
    OverUnder       PredictionType = "over_under"
    PlayerStat      PredictionType = "player_stat"
    FirstGoalScorer PredictionType = "first_goal_scorer"
)

type PredictionStatus string
const (
    Pending   PredictionStatus = "pending"
    Locked    PredictionStatus = "locked"
    Correct   PredictionStatus = "correct"
    Incorrect PredictionStatus = "incorrect"
    Void      PredictionStatus = "void"
)
```

### Validation Rules
- **Type-specific validation**: Each prediction type has unique constraints
- **Range validation**: Scores must be non-negative integers
- **Enum validation**: Outcomes must match predefined values
- **Deadline validation**: Submissions must be before event start
- **Duplicate prevention**: One prediction per user per event per type

## 2. Sports Prediction Types

### Match Outcomes
```go
type MatchOutcome struct {
    Result string `json:"result"` // "home_win", "away_win", "draw"
}
```

### Exact Scores
```go
type ExactScore struct {
    HomeScore int `json:"home_score"`
    AwayScore int `json:"away_score"`
}
```

### Over/Under Predictions
```go
type OverUnder struct {
    Line     float64 `json:"line"`     // 2.5 goals
    Choice   string  `json:"choice"`   // "over" or "under"
    Category string  `json:"category"` // "goals", "points", "corners"
}
```

### Player Statistics
```go
type PlayerStat struct {
    PlayerID   string  `json:"player_id"`
    StatType   string  `json:"stat_type"` // "goals", "assists", "points"
    Prediction float64 `json:"prediction"`
    Operator   string  `json:"operator"`  // "over", "under", "exact"
}
```

## 3. Prediction Deadlines and Time-Based Validation

### Deadline Management
```go
type EventDeadline struct {
    EventID        string    `json:"event_id"`
    StartTime      time.Time `json:"start_time"`
    DeadlineOffset int       `json:"deadline_offset_minutes"` // Minutes before start
    ActualDeadline time.Time `json:"actual_deadline"`
    IsLocked       bool      `json:"is_locked"`
}
```

### Validation Logic
- **Pre-deadline**: Accept all valid predictions
- **Grace period**: 5-15 minutes after deadline for technical issues
- **Event started**: Reject all new predictions
- **Live predictions**: Special handling for in-play events

### Time Zone Handling
- Store all times in UTC
- Convert to user's timezone for display
- Handle daylight saving transitions
- Validate against official event schedules

## 4. Scoring Algorithms and Point Calculation

### Basic Scoring System
```go
type ScoringRule struct {
    PredictionType PredictionType `json:"prediction_type"`
    CorrectPoints  int           `json:"correct_points"`
    PartialPoints  int           `json:"partial_points,omitempty"`
    BonusPoints    int           `json:"bonus_points,omitempty"`
}

// Standard scoring examples
var DefaultScoring = map[PredictionType]ScoringRule{
    MatchOutcome:    {CorrectPoints: 10},
    ExactScore:      {CorrectPoints: 40, PartialPoints: 15}, // 15 for correct outcome
    OverUnder:       {CorrectPoints: 15},
    PlayerStat:      {CorrectPoints: 20},
    FirstGoalScorer: {CorrectPoints: 25},
}
```

### Advanced Scoring Methods

#### Confidence-Based Scoring
```go
func CalculateConfidenceScore(prediction Prediction, isCorrect bool) int {
    basePoints := GetBasePoints(prediction.Type)
    if !isCorrect {
        return 0
    }
    
    confidenceMultiplier := prediction.Confidence / 100.0
    return int(float64(basePoints) * (0.5 + 0.5*confidenceMultiplier))
}
```

#### Difficulty-Based Scoring
```go
type EventDifficulty struct {
    EventID    string  `json:"event_id"`
    Difficulty float64 `json:"difficulty"` // 1.0 = normal, 2.0 = double points
    Factors    []string `json:"factors"`   // "upset", "high_scoring", "derby"
}
```

#### Progressive Scoring
- **Streak bonuses**: Consecutive correct predictions
- **Category mastery**: Bonus for expertise in specific sports
- **Early bird**: Extra points for predictions made early

## 5. Database Design for Predictions and Results

### Core Tables

#### Events Table
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    sport_type VARCHAR(50) NOT NULL,
    home_team_id UUID NOT NULL,
    away_team_id UUID NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled',
    venue VARCHAR(255),
    competition_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_status ON events(status);
```

#### Predictions Table
```sql
CREATE TABLE predictions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    contest_id UUID NOT NULL,
    event_id UUID NOT NULL,
    prediction_type VARCHAR(50) NOT NULL,
    prediction_data JSONB NOT NULL,
    confidence DECIMAL(5,2),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'pending',
    points_earned INTEGER DEFAULT 0,
    
    UNIQUE(user_id, contest_id, event_id, prediction_type)
);

CREATE INDEX idx_predictions_user_contest ON predictions(user_id, contest_id);
CREATE INDEX idx_predictions_event ON predictions(event_id);
CREATE INDEX idx_predictions_status ON predictions(status);
```

#### Results Table
```sql
CREATE TABLE event_results (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    final_score JSONB NOT NULL,
    match_outcome VARCHAR(20),
    player_stats JSONB,
    additional_stats JSONB,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Scoring Rules Table
```sql
CREATE TABLE scoring_rules (
    id UUID PRIMARY KEY,
    contest_id UUID NOT NULL,
    prediction_type VARCHAR(50) NOT NULL,
    correct_points INTEGER NOT NULL,
    partial_points INTEGER DEFAULT 0,
    bonus_multiplier DECIMAL(3,2) DEFAULT 1.0,
    is_active BOOLEAN DEFAULT true
);
```

### Optimized Queries

#### Leaderboard Query
```sql
SELECT 
    u.username,
    SUM(p.points_earned) as total_points,
    COUNT(CASE WHEN p.status = 'correct' THEN 1 END) as correct_predictions,
    COUNT(p.id) as total_predictions,
    RANK() OVER (ORDER BY SUM(p.points_earned) DESC) as rank
FROM predictions p
JOIN users u ON p.user_id = u.id
WHERE p.contest_id = $1 AND p.status IN ('correct', 'incorrect')
GROUP BY u.id, u.username
ORDER BY total_points DESC;
```

#### User Performance Query
```sql
SELECT 
    prediction_type,
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'correct' THEN 1 END) as correct,
    AVG(points_earned) as avg_points,
    SUM(points_earned) as total_points
FROM predictions 
WHERE user_id = $1 AND contest_id = $2
GROUP BY prediction_type;
```

## Best Practices Summary

### Data Integrity
- Use database constraints for referential integrity
- Implement soft deletes for audit trails
- Store prediction snapshots before events start
- Validate all inputs at API and database levels

### Performance Optimization
- Index frequently queried columns
- Use JSONB for flexible prediction data
- Implement caching for leaderboards
- Batch process scoring calculations

### Scalability Considerations
- Partition large tables by time periods
- Use read replicas for reporting queries
- Implement event sourcing for audit trails
- Consider sharding by contest or user groups

### Security & Fair Play
- Encrypt sensitive prediction data
- Implement rate limiting on submissions
- Log all prediction changes
- Detect and prevent prediction manipulation
- Use cryptographic hashes for prediction integrity