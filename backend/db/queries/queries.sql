-- name : CreateUser : exec
INSERT INTO users(username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at;

--name: GetUserByEmail: one
SELECT id, username, email, password_hash, created_at, last_login, is_active
FROM users
WHERE email=$1 AND is_active= TRUE;

--name: GetUserByID :one
SELECT id, username, email, created_at, last_login, is_active 
FROM users 
WHERE id=$1 and is_active=TRUE;

-- name: CreateGameSession :exec
INSERT INTO game_sessions (
    user_id, game_type, score, level, max_streak, total_attempts, correct_attempts, session_data, is_guest, session_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

--name: GetGameSessionByID :one 
SELECT * FROM game_sessions
WHERE id = $1;

--name: GetUserGameSession : many
SELECT * FROM game_sessions
WHERE user_id=$1 and is_guest=FALSE
ORDER BY completed_at DESC
LIMIT $2 OFFSET $3;

--name: GetGlobalLeaderboard :many
SELECT u.username, g.game_type,
        MAX(g.score) as high_score,
        COUNT(g.id) as game_played,
        AVG(g.score) as avg_score
FROM game_sessions g 
JOIN users u ON g.user_id = u.id
WHERE g.user_id IS NOT NULL AND g.is_guest=FALSE
GROUP BY u.username, g.game_type
ORDER BY high_score DESC
LIMIT $1;

-- name: UpdateGameSession :exec
UPDATE game_sessions
SET 
        score = $2,
        level =$3,
        max_streak = $4, 
        total_attempts = $5
        correct_attempts = $6, 
        session_data = $7, 
        completed_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetGuestSessionStats :one
SELECT  
        COUNT(*) as total_games,
        AVG(score) as avg_score,
        MAX(score) as high_score
FROM game_sessions
WHERE session_id = $1 and is_guest=TRUE;

-- name: TransferGuestSeesionToUser : exec
UPDATE game_sessions
SET user_id=$1, is_guest=FALSE
WHERE session_id=$2 AND user_id IS NULL;