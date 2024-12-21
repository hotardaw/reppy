-- name: CreateExerciseMuscle :one
INSERT INTO exercise_muscles (
  exercise_id,
  muscle_id,
  involvement_level
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetExerciseMuscles :many
SELECT 
  em.*,
  e.exercise_name,
  m.muscle_name,
  m.muscle_group
FROM exercise_muscles em
JOIN exercises e ON em.exercise_id = e.exercise_id
JOIN muscles m ON em.muscle_id = m.muscle_id
WHERE em.exercise_id = $1;

-- name: GetMuscleExercises :many
SELECT 
  em.*,
  e.exercise_name,
  m.muscle_name,
  m.muscle_group
FROM exercise_muscles em
JOIN exercises e ON em.exercise_id = e.exercise_id
JOIN muscles m ON em.muscle_id = m.muscle_id
WHERE em.muscle_id = $1;

-- name: UpdateExerciseMuscle :one
UPDATE exercise_muscles 
SET involvement_level = $3
WHERE exercise_id = $1 AND muscle_id = $2
RETURNING *;

-- name: DeleteExerciseMuscle :exec
DELETE FROM exercise_muscles 
WHERE exercise_id = $1 AND muscle_id = $2;

-- name: DeleteAllExerciseMuscles :exec
DELETE FROM exercise_muscles;

-- name: ListExerciseMuscles :many
SELECT 
  em.*,
  e.exercise_name,
  m.muscle_name,
  m.muscle_group
FROM exercise_muscles em
JOIN exercises e ON em.exercise_id = e.exercise_id
JOIN muscles m ON em.muscle_id = m.muscle_id
ORDER BY e.exercise_name, m.muscle_name;

-- name: ExerciseMuscleExists :one
SELECT EXISTS(
  SELECT 1 FROM exercise_muscles 
  WHERE exercise_id = $1 AND muscle_id = $2
);

-- name: GetExerciseMusclesByMuscleGroup :many
SELECT 
  em.*,
  e.exercise_name,
  m.muscle_name,
  m.muscle_group
FROM exercise_muscles em
JOIN exercises e ON em.exercise_id = e.exercise_id
JOIN muscles m ON em.muscle_id = m.muscle_id
WHERE m.muscle_group = $1
ORDER BY e.exercise_name;

-- name: GetPrimaryMusclesForExercise :many
SELECT 
  m.muscle_name,
  m.muscle_group
FROM exercise_muscles em
JOIN muscles m ON em.muscle_id = m.muscle_id
WHERE em.exercise_id = $1 
AND em.involvement_level = 'Primary'
ORDER BY m.muscle_name;