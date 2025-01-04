-- name: CreateExercise :one
INSERT INTO exercises (
    exercise_name,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetExerciseById :one
SELECT * FROM exercises 
WHERE exercise_id = $1;

-- name: GetExerciseByName :one
SELECT * FROM exercises 
WHERE exercise_name = $1;

-- name: GetAllExercises :many
SELECT * FROM exercises 
ORDER BY exercise_id;

-- name: UpdateExercise :one
UPDATE exercises 
SET exercise_name = $2, description = $3
WHERE exercise_id = $1
RETURNING *;

-- name: ExerciseExists :one
SELECT EXISTS(
  SELECT 1 FROM exercises 
  WHERE exercise_name = $1
);

-- name: SearchExercises :many
SELECT * FROM exercises 
WHERE exercise_name ILIKE $1 
ORDER BY exercise_name 
LIMIT $2;

-- name: DeleteExercise :one
DELETE FROM exercises 
WHERE exercise_id = $1
RETURNING *;

-- name: DeleteAllExercises :exec
DELETE FROM exercises;