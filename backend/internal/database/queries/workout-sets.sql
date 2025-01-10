-- RE: the "NULLIF" lines:
-- they accept text arrays, convert empty strings to NULL, then cast non-NULL values to decimal
-- name: CreateWorkoutSets :many
WITH input_rows AS (
  SELECT 
    $1::int as workout_id,
    $2::int as exercise_id,
    unnest($3::int[]) set_number,
    unnest($4::int[]) reps,
    NULLIF(unnest($5::text[]), '')::decimal resistance_value,
    NULLIF(unnest($6::text[]), '')::resistance_type_enum resistance_type,
    unnest($7::text[]) resistance_detail,
    NULLIF(unnest($8::text[]), '')::decimal rpe,
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

-- Make batch version of this later
-- name: UpdateWorkoutSetByID :one
UPDATE workout_sets 
SET 
  reps = $3,
  resistance_value = $4,
  resistance_type = $5,
  resistance_detail = $6,
  rpe = $7,
  notes = $8
WHERE workout_id = $9 
AND exercise_id = $1 
AND set_number = $2
RETURNING *;

-- name: DeleteWorkoutSetByID :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2 
AND set_number = $3;

-- name: DeleteWorkoutSetsByExercise :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2;

-- name: DeleteAllWorkoutSets :exec
DELETE FROM workout_sets
WHERE workout_id = $1;

-- Used in seeder cleanup ONLY
-- name: DeleteAllWorkoutSetsUnconditional :exec
DELETE FROM workout_sets;