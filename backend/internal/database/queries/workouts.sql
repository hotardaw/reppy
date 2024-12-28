-- CREATE: Insert a new workout
-- name: CreateWorkout :one
INSERT INTO workouts (user_id, workout_date, title)
VALUES ($1, $2, $3)
RETURNING workout_id, user_id, workout_date, title, created_at;

-- READ: Get all workouts for a specific user
-- name: GetAllWorkoutsForUser :many
SELECT workout_id, workout_date, title, created_at
FROM workouts
WHERE user_id = $1
ORDER BY workout_date DESC;

-- READ: Get a specific workout by ID
-- name: GetWorkoutByID :one
SELECT workout_id, user_id, workout_date, title, created_at
FROM workouts
WHERE workout_id = $1;

-- READ: Get workout from a date (client-side) & userID (from context)
-- name: GetWorkoutByUserIDAndDate :one
SELECT workout_id, workout_date, title, created_at
FROM workouts
WHERE user_id = $1 AND workout_date = $2;

-- UPDATE: Modify an existing workout
-- name: UpdateWorkout :one
UPDATE workouts
SET workout_date = $1,
  title = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE workout_id = $3 
AND user_id = $4
RETURNING workout_id, workout_date, title, updated_at;

-- DELETE: Remove a workout
-- name: DeleteWorkout :one
DELETE FROM workouts
WHERE workout_id = $1 
AND user_id = $2
RETURNING workout_id;