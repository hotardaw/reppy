-- name: CreateWorkoutSets :many
WITH input_rows AS (
  SELECT 
    unnest($1::int[]) workout_id,
    unnest($2::int[]) exercise_id,
    unnest($3::int[]) set_number,
    unnest($4::int[]) reps,
    unnest($5::numeric[]) resistance_value,
    unnest($6::resistance_type_enum[]) resistance_type,
    unnest($7::text[]) resistance_detail,
    unnest($8::numeric[]) rpe,
    unnest($9::text[]) notes
)
INSERT INTO workout_sets 
(workout_id, exercise_id, set_number, reps, resistance_value, resistance_type, resistance_detail, rpe, notes)
SELECT * FROM input_rows
RETURNING *;

-- name: GetAllWorkoutSets :many
SELECT 
  ws.*,
  e.exercise_name
FROM workout_sets ws
JOIN exercises e ON ws.exercise_id = e.exercise_id
WHERE ws.workout_id = $1
ORDER BY ws.exercise_id, ws.set_number;

-- name: GetAllWorkoutSetsForUserOnDate :many
SELECT 
  ws.*,
  e.exercise_name,
  w.workout_date
FROM workout_sets ws
JOIN workouts w ON ws.workout_id = w.workout_id
JOIN exercises e ON ws.exercise_id = e.exercise_id
WHERE w.user_id = $1 
AND w.workout_date = $2
ORDER BY ws.created_at, ws.exercise_id, ws.set_number;

-- Gotta make a batch version of this later
-- name: UpdateWorkoutSetDetails :one
UPDATE workout_sets 
SET 
  reps = $1,
  resistance_value = $2,
  rpe = $3,
  notes = $4
WHERE workout_id = $5 
AND exercise_id = $6 
AND set_number = $7
RETURNING *;

-- name: DeleteWorkoutSet :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2 
AND set_number = $3;

-- name: DeleteWorkoutExercise :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2;

-- name: DeleteAllWorkoutSets :exec
DELETE FROM workout_sets;