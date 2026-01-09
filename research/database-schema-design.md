# Database Schema Design for Sports Prediction Platform

## Schema Overview

The database design follows a microservices architecture with clear separation of concerns:

- **User Service**: User management and authentication
- **Contest Service**: Contest creation and management  
- **Prediction Service**: User predictions and validation
- **Scoring Service**: Points calculation and leaderboards
- **Sports Service**: Events, teams, and results

## Core Entity Relationships

```
Users (1) -----> (M) Predictions
Contests (1) ---> (M) Predictions  
Events (1) -----> (M) Predictions
Events (1) -----> (1) EventResults
Contests (1) ---> (M) ScoringRules
```

## Detailed Schema

### Users Service

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    avatar_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'en',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User preferences
CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    notification_settings JSONB DEFAULT '{}',
    privacy_settings JSONB DEFAULT '{}',
    display_preferences JSONB DEFAULT '{}'
);
```

### Sports Service

```sql
-- Sports and competitions
CREATE TABLE sports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL, -- 'football', 'basketball'
    icon_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE competitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id),
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) NOT NULL,
    country VARCHAR(100),
    season VARCHAR(20),
    is_active BOOLEAN DEFAULT true
);

-- Teams
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    short_name VARCHAR(50),
    code VARCHAR(10),
    logo_url VARCHAR(500),
    sport_id UUID NOT NULL REFERENCES sports(id),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true
);

-- Players
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    position VARCHAR(50),
    jersey_number INTEGER,
    team_id UUID REFERENCES teams(id),
    sport_id UUID NOT NULL REFERENCES sports(id),
    is_active BOOLEAN DEFAULT true
);

-- Events (matches/games)
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    competition_id UUID NOT NULL REFERENCES competitions(id),
    home_team_id UUID NOT NULL REFERENCES teams(id),
    away_team_id UUID NOT NULL REFERENCES teams(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    venue VARCHAR(255),
    status VARCHAR(20) DEFAULT 'scheduled', -- scheduled, live, finished, cancelled
    round VARCHAR(50),
    matchday INTEGER,
    external_id VARCHAR(100), -- For API integration
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Event results
CREATE TABLE event_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL UNIQUE REFERENCES events(id),
    home_score INTEGER NOT NULL DEFAULT 0,
    away_score INTEGER NOT NULL DEFAULT 0,
    match_outcome VARCHAR(20), -- home_win, away_win, draw
    additional_stats JSONB DEFAULT '{}', -- corners, cards, possession, etc.
    player_stats JSONB DEFAULT '{}', -- goals, assists, etc.
    confirmed BOOLEAN DEFAULT false,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Contest Service

```sql
-- Contests
CREATE TABLE contests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    creator_id UUID NOT NULL, -- References users(id)
    sport_id UUID NOT NULL REFERENCES sports(id),
    competition_id UUID REFERENCES competitions(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    max_participants INTEGER,
    entry_fee DECIMAL(10,2) DEFAULT 0,
    prize_pool DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft', -- draft, active, finished, cancelled
    is_public BOOLEAN DEFAULT true,
    rules JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Contest participants
CREATE TABLE contest_participants (
    contest_id UUID NOT NULL REFERENCES contests(id),
    user_id UUID NOT NULL, -- References users(id)
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    PRIMARY KEY (contest_id, user_id)
);

-- Scoring rules per contest
CREATE TABLE scoring_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contest_id UUID NOT NULL REFERENCES contests(id),
    prediction_type VARCHAR(50) NOT NULL,
    correct_points INTEGER NOT NULL DEFAULT 0,
    partial_points INTEGER DEFAULT 0,
    bonus_multiplier DECIMAL(3,2) DEFAULT 1.0,
    difficulty_multiplier DECIMAL(3,2) DEFAULT 1.0,
    is_active BOOLEAN DEFAULT true
);

-- Contest events (which events are included)
CREATE TABLE contest_events (
    contest_id UUID NOT NULL REFERENCES contests(id),
    event_id UUID NOT NULL REFERENCES events(id),
    deadline_offset INTEGER DEFAULT 0, -- Minutes before event start
    is_active BOOLEAN DEFAULT true,
    PRIMARY KEY (contest_id, event_id)
);
```

### Prediction Service

```sql
-- Predictions
CREATE TABLE predictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- References users(id)
    contest_id UUID NOT NULL REFERENCES contests(id),
    event_id UUID NOT NULL REFERENCES events(id),
    prediction_type VARCHAR(50) NOT NULL,
    prediction_data JSONB NOT NULL,
    confidence DECIMAL(5,2) CHECK (confidence >= 0 AND confidence <= 100),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    locked_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'pending', -- pending, locked, correct, incorrect, void
    points_earned INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, contest_id, event_id, prediction_type)
);

-- Prediction validation rules
CREATE TABLE prediction_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id),
    prediction_type VARCHAR(50) NOT NULL,
    validation_schema JSONB NOT NULL,
    display_config JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true
);
```

### Scoring Service

```sql
-- User contest scores (aggregated)
CREATE TABLE user_contest_scores (
    user_id UUID NOT NULL, -- References users(id)
    contest_id UUID NOT NULL REFERENCES contests(id),
    total_points INTEGER DEFAULT 0,
    correct_predictions INTEGER DEFAULT 0,
    total_predictions INTEGER DEFAULT 0,
    accuracy_percentage DECIMAL(5,2) DEFAULT 0,
    rank INTEGER,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, contest_id)
);

-- Leaderboard snapshots for performance
CREATE TABLE leaderboard_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contest_id UUID NOT NULL REFERENCES contests(id),
    snapshot_data JSONB NOT NULL,
    snapshot_type VARCHAR(20) NOT NULL, -- daily, weekly, final
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User achievements
CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- References users(id)
    achievement_type VARCHAR(50) NOT NULL,
    achievement_data JSONB DEFAULT '{}',
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    contest_id UUID REFERENCES contests(id)
);
```

## Indexes for Performance

```sql
-- Events
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_competition ON events(competition_id);

-- Predictions
CREATE INDEX idx_predictions_user_contest ON predictions(user_id, contest_id);
CREATE INDEX idx_predictions_event ON predictions(event_id);
CREATE INDEX idx_predictions_status ON predictions(status);
CREATE INDEX idx_predictions_submitted_at ON predictions(submitted_at);

-- Contest participants
CREATE INDEX idx_contest_participants_user ON contest_participants(user_id);

-- User contest scores
CREATE INDEX idx_user_scores_contest ON user_contest_scores(contest_id, total_points DESC);
CREATE INDEX idx_user_scores_rank ON user_contest_scores(contest_id, rank);
```

## Data Types and Validation

### Prediction Data Examples

```json
// Match Outcome
{
  "result": "home_win" // home_win, away_win, draw
}

// Exact Score
{
  "home_score": 2,
  "away_score": 1
}

// Over/Under
{
  "line": 2.5,
  "choice": "over", // over, under
  "category": "goals" // goals, corners, cards
}

// Player Stats
{
  "player_id": "uuid",
  "stat_type": "goals", // goals, assists, points, saves
  "prediction": 1,
  "operator": "over" // over, under, exact
}

// First Goal Scorer
{
  "player_id": "uuid",
  "no_goal": false // true if predicting no goals
}
```

## Partitioning Strategy

For high-volume data, consider partitioning:

```sql
-- Partition predictions by month
CREATE TABLE predictions_y2024m01 PARTITION OF predictions
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Partition events by year
CREATE TABLE events_2024 PARTITION OF events
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
```

## Data Retention Policy

```sql
-- Archive old predictions (keep for 2 years)
DELETE FROM predictions 
WHERE created_at < NOW() - INTERVAL '2 years';

-- Archive old events (keep for 5 years)
DELETE FROM events 
WHERE start_time < NOW() - INTERVAL '5 years';
```