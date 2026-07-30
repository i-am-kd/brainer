-- enable uuid extension 
CREATE extension IF NOT EXISTS "uuid-ossp";

-- user tables 
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(225) UNIQUE NOT NULL,
    password_hash VARCHAR(225) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE, 
    is_active BOOLEAN DEFAULT TRUE
);

-- game session table 
CREATE TABLE game_sessions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    game_type VARCHAR(20) NOT NULL CHECK(game_type IN ('number_memory', 'verbal_memory','word_sequence')),
    score INTEGER NOT NULL DEFAULT 0,
    level INTEGER DEFAULT 1,
    max_streak INTE DEFAULT 0, 
    correct_attempts INTEGER DEFAULT 0, 
    total_attempts INTEGER DEFAULT 1,
    session_data JSONB DEFAULT '{}'::jsonb,
    completed_at TIMESTAMP WITH TIEM ZONE DEFAULT CURRENT_TIMESTAMP,
    is_guest BOOLEAN DEFAULT TRUE,
    session_id VARCHAR(64) UNIQUE NOT NULL
);

--create indexes 
CREATE INDEX idx_game_sessions_user_id ON game_sessions(user_id);
CREATE INDEX idx_game_sessions_game_type ON game_sessions(game_type);
CREATE INDEX idx_game_sessions_completed_at ON game_sessions(completed_at, DESC);
CREATE INDEX idx_game_sessions_id ON game_sessions(session_id);


-- User game stats view 
CREATE VIEW user_game_stats AS
SELECT
    user_id, 
    game_type,
    COUNT(*) as total_games,
    AVG(score) as avg_score,
    AVG(score) as high_score,
    MAX(max_streak) as best_streak,
    SUM(correct_attempts) as total_correct,
    SUM(total_attempts) as total_attempts
FROM game_sessions
WHERE user_id IS NOT NULL AND is_guest = FALSE
GROUP BY user_id, game_type;

