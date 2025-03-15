CREATE OR REPLACE FUNCTION now() RETURNS TIMESTAMPTZ AS $$
  SELECT CURRENT_TIMESTAMP;
$$ LANGUAGE SQL;

-- Users

-- name: CreateUser :execresult
INSERT INTO users (id, created_at, data)
VALUES ($1, $2, $3)
RETURNING id, created_at, data;

-- name: UpdateUser :execresult
UPDATE users
SET updated_at = $2, data = $3
WHERE id = $1
RETURNING id, created_at, updated_at, data;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE id = $1
RETURNING id, created_at, updated_at, data;

-- name: GetUser :one
SELECT id, created_at, updated_at, data
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, created_at, updated_at, data
FROM users
ORDER BY created_at DESC;

-- Secrets

-- name: CreateSecret :execresult
INSERT INTO secrets (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: UpdateSecret :execresult
UPDATE secrets
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: DeleteSecret :execresult
DELETE FROM secrets
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: GetSecret :one
SELECT id, user_id, created_at, updated_at, data
FROM secrets
WHERE id = $1 AND user_id = $2;

-- name: GetAllSecrets :many
SELECT id, user_id, created_at, updated_at, data
FROM secrets
ORDER BY created_at DESC;

-- API Tokens

-- name: CreateAPIToken :execresult
INSERT INTO api_tokens (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: UpdateAPIToken :execresult
UPDATE api_tokens
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: DeleteAPIToken :execresult
DELETE FROM api_tokens
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: GetAPIToken :one
SELECT id, user_id, created_at, updated_at, data
FROM api_tokens
WHERE id = $1 AND user_id = $2;

-- name: GetAllAPITokens :many
SELECT id, user_id, created_at, updated_at, data
FROM api_tokens
ORDER BY created_at DESC;

-- Rate Limits

-- name: CreateRateLimit :execresult
INSERT INTO rate_limits (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: UpdateRateLimit :execresult
UPDATE rate_limits
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: DeleteRateLimit :execresult
DELETE FROM rate_limits
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: GetRateLimit :one
SELECT id, user_id, created_at, updated_at, data
FROM rate_limits
WHERE id = $1 AND user_id = $2;

-- name: GetAllRateLimits :many
SELECT id, user_id, created_at, updated_at, data
FROM rate_limits
ORDER BY created_at DESC;

-- Chat Sessions

-- name: CreateChatSession :execresult
INSERT INTO chat_sessions (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: UpdateChatSession :execresult
UPDATE chat_sessions
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: DeleteChatSession :execresult
DELETE FROM chat_sessions
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: GetChatSession :one
SELECT id, user_id, created_at, updated_at, data
FROM chat_sessions
WHERE id = $1 AND user_id = $2;

-- name: GetAllChatSessions :many
SELECT id, user_id, created_at, updated_at, data
FROM chat_sessions
ORDER BY created_at DESC;

-- name: GetChatSessionsByUserID :many
SELECT id, user_id, created_at, updated_at, data
FROM chat_sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- Chat Session Items

-- name: CreateChatSessionItem :execresult
INSERT INTO chat_session_items (id, user_id, created_at, type, session_id, data)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, session_id, created_at, updated_at, type, data;

-- name: UpdateChatSessionItem :execresult
UPDATE chat_session_items
SET updated_at = $3, type = $4, session_id = $5, data = $6
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, session_id, created_at, updated_at, type, data;

-- name: DeleteChatSessionItem :execresult
DELETE FROM chat_session_items
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, session_id, created_at, updated_at, type, data;

-- name: DeleteChatSessionItemsBySessionID :exec
DELETE FROM chat_session_items
WHERE session_id = $1 AND user_id = $2;

-- name: GetChatSessionItem :one
SELECT id, user_id, session_id, created_at, updated_at, type, data
FROM chat_session_items
WHERE id = $1 AND user_id = $2;

-- name: GetAllChatSessionItems :many
SELECT id, user_id, session_id, created_at, updated_at, type, data
FROM chat_session_items
ORDER BY created_at DESC;

-- name: GetChatSessionItemsBySessionID :many
SELECT id, user_id, session_id, created_at, updated_at, type, data
FROM chat_session_items
WHERE session_id = $1 AND user_id = $2
ORDER BY created_at DESC;

-- Game States

-- name: CreateGameState :execresult
INSERT INTO game_states (id, created_at, data)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at, data;

-- name: UpdateGameState :execresult
UPDATE game_states
SET updated_at = $2, data = $3
WHERE id = $1
RETURNING id, created_at, updated_at, data;

-- name: DeleteGameState :execresult
DELETE FROM game_states
WHERE id = $1
RETURNING id, created_at, updated_at, data;

-- name: GetGameState :one
SELECT id, created_at, updated_at, data
FROM game_states
WHERE id = $1;

-- name: GetAllGameStates :many
SELECT id, created_at, updated_at, data
FROM game_states
ORDER BY created_at DESC;

-- Game State Versions

-- name: CreateGameStateVersion :execresult
INSERT INTO game_state_versions (id, game_state_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, game_state_id, created_at, updated_at, data;

-- name: GetGameStateVersion :one
SELECT id, game_state_id, created_at, updated_at, data
FROM game_state_versions
WHERE id = $1;

-- name: GetAllGameStateVersions :many
SELECT id, game_state_id, created_at, updated_at, data
FROM game_state_versions
ORDER BY created_at DESC;

-- name: GetGameStateVersionsByGameStateID :many
SELECT id, game_state_id, created_at, updated_at, data
FROM game_state_versions
WHERE game_state_id = $1
ORDER BY created_at DESC;

-- name: GetLatestGameStateVersionByGameStateID :one
SELECT id, game_state_id, created_at, updated_at, data
FROM game_state_versions
WHERE game_state_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- Cards

-- name: CreateCard :execresult
INSERT INTO cards (id, type, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, type, data;

-- name: UpdateCard :execresult
UPDATE cards
SET updated_at = $2, type = $3, data = $4
WHERE id = $1
RETURNING id, created_at, updated_at, type, data;

-- name: DeleteCard :execresult
DELETE FROM cards
WHERE id = $1
RETURNING id, created_at, updated_at, type, data;

-- name: GetCard :one
SELECT id, created_at, updated_at, type, data
FROM cards
WHERE id = $1;

-- name: GetAllCards :many
SELECT id, created_at, updated_at, type, data
FROM cards
ORDER BY created_at DESC;

-- Decks

-- name: CreateDeck :execresult
INSERT INTO decks (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: UpdateDeck :execresult
UPDATE decks
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: DeleteDeck :execresult
DELETE FROM decks
WHERE id = $1 AND user_id = $2;

-- name: GetDeck :one
SELECT id, user_id, created_at, updated_at, data
FROM decks
WHERE id = $1 AND user_id = $2;

-- name: GetAllDecks :many
SELECT id, user_id, created_at, updated_at, data
FROM decks
ORDER BY created_at DESC;


-- Integrations

-- name: CreateIntegration :execresult
INSERT INTO integrations (id, type, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, type, created_at, data;

-- name: UpdateIntegration :execresult
UPDATE integrations
SET updated_at = $2, data = $3
WHERE id = $1
RETURNING id, type, created_at, updated_at, data;

-- name: DeleteIntegration :execresult
DELETE FROM integrations
WHERE id = $1
RETURNING id, type, created_at, updated_at, data;

-- name: GetIntegration :one
SELECT id, type, created_at, updated_at, data
FROM integrations
WHERE id = $1;

-- name: GetAllIntegrations :many
SELECT id, type, created_at, updated_at, data
FROM integrations
ORDER BY created_at DESC;
