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

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, username)
VALUES ($1, $2, $3)
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

-- name: DeleteUser :exec
UPDATE users
SET active = false, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;