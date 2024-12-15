-- The comments above each query are SQLc directives dictating the naming of the generated Go func & what type of query it is (:one, :many, or :exec).
--  :many returns a slice of records via QueryContext
-- :one returns a single record via QueryRowContext
-- :exec returns the error from ExecContext
-- More: https://docs.sqlc.dev/en/latest/reference/query-annotations.html

-- name: GetUser :one
SELECT users.*
FROM users
WHERE users.user_id = $1;

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, username)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, username = $4
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1
RETURNING *;