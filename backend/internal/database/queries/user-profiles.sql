-- name: GetUserProfile :one
SELECT user_profiles.*
FROM user_profiles
WHERE user_profiles.user_id = $1;

-- name: GetAllUserProfiles :many
SELECT user_profiles.*
FROM user_profiles;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, first_name, last_name, date_of_birth, gender, height_inches)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET 
  first_name = $2, 
  last_name = $3, 
  height_inches = $4, 
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: DeleteUserProfile :one
DELETE FROM user_profiles
WHERE user_id = $1
RETURNING *;

-- name: DeleteAllUserProfiles :exec
DELETE FROM user_profiles;