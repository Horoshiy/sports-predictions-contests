-- Initialize database for sports prediction contests platform

-- Create users table (will be auto-created by GORM, but this ensures it exists)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create scores table for scoring service
CREATE TABLE IF NOT EXISTS scores (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    contest_id INTEGER NOT NULL,
    prediction_id INTEGER NOT NULL,
    points DECIMAL(10,2) NOT NULL DEFAULT 0,
    scored_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, contest_id, prediction_id)
);

-- Create indexes for scores table
CREATE INDEX IF NOT EXISTS idx_scores_contest_id ON scores(contest_id);
CREATE INDEX IF NOT EXISTS idx_scores_user_id ON scores(user_id);
CREATE INDEX IF NOT EXISTS idx_scores_prediction_id ON scores(prediction_id);
CREATE INDEX IF NOT EXISTS idx_scores_deleted_at ON scores(deleted_at);

-- Create leaderboards table for ranking
CREATE TABLE IF NOT EXISTS leaderboards (
    id SERIAL PRIMARY KEY,
    contest_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    total_points DECIMAL(10,2) NOT NULL DEFAULT 0,
    rank INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(contest_id, user_id)
);

-- Create indexes for leaderboards table
CREATE INDEX IF NOT EXISTS idx_leaderboards_contest_id ON leaderboards(contest_id);
CREATE INDEX IF NOT EXISTS idx_leaderboards_contest_rank ON leaderboards(contest_id, rank);
CREATE INDEX IF NOT EXISTS idx_leaderboards_deleted_at ON leaderboards(deleted_at);

-- Create user_streaks table for streak tracking
CREATE TABLE IF NOT EXISTS user_streaks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    contest_id INTEGER NOT NULL,
    current_streak INTEGER NOT NULL DEFAULT 0,
    max_streak INTEGER NOT NULL DEFAULT 0,
    last_prediction_id INTEGER,
    last_prediction_correct BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, contest_id)
);

CREATE INDEX IF NOT EXISTS idx_user_streaks_contest_user ON user_streaks(contest_id, user_id);
CREATE INDEX IF NOT EXISTS idx_user_streaks_deleted_at ON user_streaks(deleted_at);

-- Create sports table
CREATE TABLE IF NOT EXISTS sports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    external_id VARCHAR(50) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_sports_slug ON sports(slug);
CREATE INDEX IF NOT EXISTS idx_sports_is_active ON sports(is_active);
CREATE INDEX IF NOT EXISTS idx_sports_external_id ON sports(external_id);
CREATE INDEX IF NOT EXISTS idx_sports_deleted_at ON sports(deleted_at);

-- Create leagues table
CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    country VARCHAR(100),
    season VARCHAR(50),
    external_id VARCHAR(50) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_leagues_sport_id ON leagues(sport_id);
CREATE INDEX IF NOT EXISTS idx_leagues_slug ON leagues(slug);
CREATE INDEX IF NOT EXISTS idx_leagues_is_active ON leagues(is_active);
CREATE INDEX IF NOT EXISTS idx_leagues_external_id ON leagues(external_id);
CREATE INDEX IF NOT EXISTS idx_leagues_deleted_at ON leagues(deleted_at);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    short_name VARCHAR(50),
    logo_url VARCHAR(500),
    country VARCHAR(100),
    external_id VARCHAR(50) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_teams_sport_id ON teams(sport_id);
CREATE INDEX IF NOT EXISTS idx_teams_slug ON teams(slug);
CREATE INDEX IF NOT EXISTS idx_teams_is_active ON teams(is_active);
CREATE INDEX IF NOT EXISTS idx_teams_external_id ON teams(external_id);
CREATE INDEX IF NOT EXISTS idx_teams_deleted_at ON teams(deleted_at);

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    league_id INTEGER NOT NULL REFERENCES leagues(id) ON DELETE RESTRICT,
    home_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    away_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    home_score INTEGER DEFAULT 0,
    away_score INTEGER DEFAULT 0,
    result_data TEXT,
    external_id VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_matches_league_id ON matches(league_id);
CREATE INDEX IF NOT EXISTS idx_matches_home_team_id ON matches(home_team_id);
CREATE INDEX IF NOT EXISTS idx_matches_away_team_id ON matches(away_team_id);
CREATE INDEX IF NOT EXISTS idx_matches_scheduled_at ON matches(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);
CREATE INDEX IF NOT EXISTS idx_matches_external_id ON matches(external_id);
CREATE INDEX IF NOT EXISTS idx_matches_deleted_at ON matches(deleted_at);


-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data TEXT,
    channel VARCHAR(20) NOT NULL DEFAULT 'in_app',
    is_read BOOLEAN DEFAULT false,
    sent_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

-- Create notification_preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    channel VARCHAR(20) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    telegram_chat_id BIGINT,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, channel)
);

CREATE INDEX IF NOT EXISTS idx_notification_preferences_user ON notification_preferences(user_id);


-- Create user_teams table (for team tournaments feature)
CREATE TABLE IF NOT EXISTS user_teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    invite_code VARCHAR(20) UNIQUE NOT NULL,
    captain_id INTEGER NOT NULL,
    max_members INTEGER DEFAULT 10,
    current_members INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_user_teams_invite_code ON user_teams(invite_code);
CREATE INDEX IF NOT EXISTS idx_user_teams_captain_id ON user_teams(captain_id);
CREATE INDEX IF NOT EXISTS idx_user_teams_is_active ON user_teams(is_active);
CREATE INDEX IF NOT EXISTS idx_user_teams_deleted_at ON user_teams(deleted_at);

-- Create user_team_members table
CREATE TABLE IF NOT EXISTS user_team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES user_teams(id) ON DELETE RESTRICT,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(team_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_team_members_team_id ON user_team_members(team_id);
CREATE INDEX IF NOT EXISTS idx_user_team_members_user_id ON user_team_members(user_id);
CREATE INDEX IF NOT EXISTS idx_user_team_members_deleted_at ON user_team_members(deleted_at);

-- Create user_team_contest_entries table (links user teams to contests)
CREATE TABLE IF NOT EXISTS user_team_contest_entries (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES user_teams(id) ON DELETE RESTRICT,
    contest_id INTEGER NOT NULL,
    total_points DECIMAL(10,2) NOT NULL DEFAULT 0,
    rank INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(team_id, contest_id)
);

CREATE INDEX IF NOT EXISTS idx_user_team_contest_entries_team_id ON user_team_contest_entries(team_id);
CREATE INDEX IF NOT EXISTS idx_user_team_contest_entries_contest_id ON user_team_contest_entries(contest_id);
CREATE INDEX IF NOT EXISTS idx_user_team_contest_entries_rank ON user_team_contest_entries(contest_id, rank);
CREATE INDEX IF NOT EXISTS idx_user_team_contest_entries_deleted_at ON user_team_contest_entries(deleted_at);


-- Create prop_types table for props predictions
CREATE TABLE IF NOT EXISTS prop_types (
    id SERIAL PRIMARY KEY,
    sport_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    value_type VARCHAR(20) NOT NULL,
    default_line DECIMAL(10,2),
    min_value DECIMAL(10,2),
    max_value DECIMAL(10,2),
    points_correct DECIMAL(10,2) NOT NULL DEFAULT 2,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(sport_type, slug)
);

CREATE INDEX IF NOT EXISTS idx_prop_types_sport ON prop_types(sport_type);
CREATE INDEX IF NOT EXISTS idx_prop_types_category ON prop_types(category);
CREATE INDEX IF NOT EXISTS idx_prop_types_is_active ON prop_types(is_active);
CREATE INDEX IF NOT EXISTS idx_prop_types_deleted_at ON prop_types(deleted_at);

-- Insert default prop types for Soccer
INSERT INTO prop_types (sport_type, name, slug, description, category, value_type, default_line, points_correct) VALUES
('Soccer', 'Total Goals Over/Under', 'total-goals-ou', 'Predict if total goals will be over or under the line', 'match', 'over_under', 2.5, 2),
('Soccer', 'Total Corners Over/Under', 'total-corners-ou', 'Predict if total corners will be over or under the line', 'match', 'over_under', 9.5, 2),
('Soccer', 'Both Teams to Score', 'btts', 'Predict if both teams will score', 'match', 'yes_no', NULL, 2),
('Soccer', 'First Team to Score', 'first-to-score', 'Predict which team scores first', 'match', 'team_select', NULL, 3),
('Soccer', 'Player to Score Anytime', 'player-goal', 'Predict if a specific player will score', 'player', 'yes_no', NULL, 4),
('Soccer', 'Total Cards Over/Under', 'total-cards-ou', 'Predict if total cards will be over or under the line', 'match', 'over_under', 3.5, 2)
ON CONFLICT (sport_type, slug) DO NOTHING;
