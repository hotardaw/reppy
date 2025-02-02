-- The comments above each query are SQLc directives dictating the naming of the generated Go func & what type of query it is (:one, :many, or :exec).
--  :many returns a slice of records via QueryContext
-- :one returns a single record via QueryRowContext
-- :exec returns the error from ExecContext
-- More: https://docs.sqlc.dev/en/latest/reference/query-annotations.html

-- name: GetUser :one
SELECT users.*
FROM users
WHERE users.user_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 AND active = true
LIMIT 1;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1;

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, username)
VALUES ($1, $2, $3)
RETURNING *;

-- The password_hash here is a dummy hash - Google users don't need a password, and I don't feel like changing the table's constraints
-- name: CreateGoogleUser :one
INSERT INTO users (email, username, google_id, auth_provider, password_hash)
VALUES ($1, $2, $3, 'google', 'GOOGLE_AUTH_USER')
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
  email = COALESCE($2, email),
  password_hash = COALESCE($3, password_hash),
  username = COALESCE($4, username),
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- Soft delete only - too many headaches if this gets actual users.
-- name: DeleteUser :one
UPDATE users
SET active = false, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;