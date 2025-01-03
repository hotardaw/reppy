-- name: GetUserProfile :one
SELECT user_profiles.*
FROM user_profiles
WHERE user_profiles.user_id = $1;

-- name: GetAllUserProfiles :many
SELECT user_profiles.*
FROM user_profiles;

-- Weird SQLc error with using user_profiles.* here - generated empty select statement. Explicitly naming columns to select instead.
-- name: GetAllActiveUserProfiles :many
SELECT
  ups.profile_id,
  ups.user_id,
  ups.first_name,
  ups.last_name,
  ups.date_of_birth,
  ups.gender,
  ups.height_inches,
  ups.weight_pounds,
  ups.profile_picture_url,
  ups.created_at,
  ups.updated_at,
  u.active
FROM user_profiles ups
JOIN users u ON ups.user_id = u.user_id
WHERE u.active = true;

-- Same thing here - explicitly naming columns to return.
-- name: GetAllInactiveUserProfiles :many
SELECT 
  ups.profile_id,
  ups.user_id,
  ups.first_name,
  ups.last_name,
  ups.date_of_birth,
  ups.gender,
  ups.height_inches,
  ups.weight_pounds,
  ups.profile_picture_url,
  ups.created_at,
  ups.updated_at,
  u.active
FROM user_profiles ups
JOIN users u ON ups.user_id = u.user_id
WHERE u.active = false;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, first_name, last_name, date_of_birth, gender, height_inches, weight_pounds)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET 
  first_name = $2, 
  last_name = $3, 
  date_of_birth = $4,
  gender = $5,
  height_inches = $6,
  weight_pounds = $7,
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: DeleteUserProfile :one
DELETE FROM user_profiles
WHERE user_id = $1
RETURNING *;

-- name: DeleteAllUserProfiles :exec
DELETE FROM user_profiles;