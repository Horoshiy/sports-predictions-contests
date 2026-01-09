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
