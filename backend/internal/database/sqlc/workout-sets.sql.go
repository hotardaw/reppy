// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workout-sets.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createWorkoutSets = `-- name: CreateWorkoutSets :many
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
SELECT workout_id, exercise_id, set_number, reps, resistance_value, resistance_type, resistance_detail, rpe, notes FROM input_rows
RETURNING workout_id, exercise_id, set_number, reps, resistance_value, resistance_type, resistance_detail, rpe, percent_1rm, notes, created_at
`

type CreateWorkoutSetsParams struct {
	Column1 []int32
	Column2 []int32
	Column3 []int32
	Column4 []int32
	Column5 []string
	Column6 []ResistanceTypeEnum
	Column7 []string
	Column8 []string
	Column9 []string
}

func (q *Queries) CreateWorkoutSets(ctx context.Context, arg CreateWorkoutSetsParams) ([]WorkoutSet, error) {
	rows, err := q.db.QueryContext(ctx, createWorkoutSets,
		pq.Array(arg.Column1),
		pq.Array(arg.Column2),
		pq.Array(arg.Column3),
		pq.Array(arg.Column4),
		pq.Array(arg.Column5),
		pq.Array(arg.Column6),
		pq.Array(arg.Column7),
		pq.Array(arg.Column8),
		pq.Array(arg.Column9),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WorkoutSet
	for rows.Next() {
		var i WorkoutSet
		if err := rows.Scan(
			&i.WorkoutID,
			&i.ExerciseID,
			&i.SetNumber,
			&i.Reps,
			&i.ResistanceValue,
			&i.ResistanceType,
			&i.ResistanceDetail,
			&i.Rpe,
			&i.Percent1rm,
			&i.Notes,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteAllWorkoutSets = `-- name: DeleteAllWorkoutSets :exec
DELETE FROM workout_sets
`

func (q *Queries) DeleteAllWorkoutSets(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllWorkoutSets)
	return err
}

const deleteWorkoutExercise = `-- name: DeleteWorkoutExercise :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2
`

type DeleteWorkoutExerciseParams struct {
	WorkoutID  int32
	ExerciseID int32
}

func (q *Queries) DeleteWorkoutExercise(ctx context.Context, arg DeleteWorkoutExerciseParams) error {
	_, err := q.db.ExecContext(ctx, deleteWorkoutExercise, arg.WorkoutID, arg.ExerciseID)
	return err
}

const deleteWorkoutSet = `-- name: DeleteWorkoutSet :exec
DELETE FROM workout_sets 
WHERE workout_id = $1 
AND exercise_id = $2 
AND set_number = $3
`

type DeleteWorkoutSetParams struct {
	WorkoutID  int32
	ExerciseID int32
	SetNumber  int32
}

func (q *Queries) DeleteWorkoutSet(ctx context.Context, arg DeleteWorkoutSetParams) error {
	_, err := q.db.ExecContext(ctx, deleteWorkoutSet, arg.WorkoutID, arg.ExerciseID, arg.SetNumber)
	return err
}

const getAllWorkoutSets = `-- name: GetAllWorkoutSets :many
SELECT 
  ws.workout_id, ws.exercise_id, ws.set_number, ws.reps, ws.resistance_value, ws.resistance_type, ws.resistance_detail, ws.rpe, ws.percent_1rm, ws.notes, ws.created_at,
  e.exercise_name
FROM workout_sets ws
JOIN exercises e ON ws.exercise_id = e.exercise_id
WHERE ws.workout_id = $1
ORDER BY ws.exercise_id, ws.set_number
`

type GetAllWorkoutSetsRow struct {
	WorkoutID        int32
	ExerciseID       int32
	SetNumber        int32
	Reps             sql.NullInt32
	ResistanceValue  sql.NullInt32
	ResistanceType   NullResistanceTypeEnum
	ResistanceDetail sql.NullString
	Rpe              sql.NullString
	Percent1rm       sql.NullString
	Notes            sql.NullString
	CreatedAt        sql.NullTime
	ExerciseName     string
}

func (q *Queries) GetAllWorkoutSets(ctx context.Context, workoutID int32) ([]GetAllWorkoutSetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllWorkoutSets, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllWorkoutSetsRow
	for rows.Next() {
		var i GetAllWorkoutSetsRow
		if err := rows.Scan(
			&i.WorkoutID,
			&i.ExerciseID,
			&i.SetNumber,
			&i.Reps,
			&i.ResistanceValue,
			&i.ResistanceType,
			&i.ResistanceDetail,
			&i.Rpe,
			&i.Percent1rm,
			&i.Notes,
			&i.CreatedAt,
			&i.ExerciseName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllWorkoutSetsForUserOnDate = `-- name: GetAllWorkoutSetsForUserOnDate :many
SELECT 
  ws.workout_id, ws.exercise_id, ws.set_number, ws.reps, ws.resistance_value, ws.resistance_type, ws.resistance_detail, ws.rpe, ws.percent_1rm, ws.notes, ws.created_at,
  e.exercise_name,
  w.workout_date
FROM workout_sets ws
JOIN workouts w ON ws.workout_id = w.workout_id
JOIN exercises e ON ws.exercise_id = e.exercise_id
WHERE w.user_id = $1 
AND w.workout_date = $2
ORDER BY ws.created_at, ws.exercise_id, ws.set_number
`

type GetAllWorkoutSetsForUserOnDateParams struct {
	UserID      sql.NullInt32
	WorkoutDate time.Time
}

type GetAllWorkoutSetsForUserOnDateRow struct {
	WorkoutID        int32
	ExerciseID       int32
	SetNumber        int32
	Reps             sql.NullInt32
	ResistanceValue  sql.NullInt32
	ResistanceType   NullResistanceTypeEnum
	ResistanceDetail sql.NullString
	Rpe              sql.NullString
	Percent1rm       sql.NullString
	Notes            sql.NullString
	CreatedAt        sql.NullTime
	ExerciseName     string
	WorkoutDate      time.Time
}

func (q *Queries) GetAllWorkoutSetsForUserOnDate(ctx context.Context, arg GetAllWorkoutSetsForUserOnDateParams) ([]GetAllWorkoutSetsForUserOnDateRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllWorkoutSetsForUserOnDate, arg.UserID, arg.WorkoutDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllWorkoutSetsForUserOnDateRow
	for rows.Next() {
		var i GetAllWorkoutSetsForUserOnDateRow
		if err := rows.Scan(
			&i.WorkoutID,
			&i.ExerciseID,
			&i.SetNumber,
			&i.Reps,
			&i.ResistanceValue,
			&i.ResistanceType,
			&i.ResistanceDetail,
			&i.Rpe,
			&i.Percent1rm,
			&i.Notes,
			&i.CreatedAt,
			&i.ExerciseName,
			&i.WorkoutDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateWorkoutSetDetails = `-- name: UpdateWorkoutSetDetails :one
UPDATE workout_sets 
SET 
  reps = $1,
  resistance_value = $2,
  rpe = $3,
  notes = $4
WHERE workout_id = $5 
AND exercise_id = $6 
AND set_number = $7
RETURNING workout_id, exercise_id, set_number, reps, resistance_value, resistance_type, resistance_detail, rpe, percent_1rm, notes, created_at
`

type UpdateWorkoutSetDetailsParams struct {
	Reps            sql.NullInt32
	ResistanceValue sql.NullInt32
	Rpe             sql.NullString
	Notes           sql.NullString
	WorkoutID       int32
	ExerciseID      int32
	SetNumber       int32
}

// Gotta make a batch version of this later
func (q *Queries) UpdateWorkoutSetDetails(ctx context.Context, arg UpdateWorkoutSetDetailsParams) (WorkoutSet, error) {
	row := q.db.QueryRowContext(ctx, updateWorkoutSetDetails,
		arg.Reps,
		arg.ResistanceValue,
		arg.Rpe,
		arg.Notes,
		arg.WorkoutID,
		arg.ExerciseID,
		arg.SetNumber,
	)
	var i WorkoutSet
	err := row.Scan(
		&i.WorkoutID,
		&i.ExerciseID,
		&i.SetNumber,
		&i.Reps,
		&i.ResistanceValue,
		&i.ResistanceType,
		&i.ResistanceDetail,
		&i.Rpe,
		&i.Percent1rm,
		&i.Notes,
		&i.CreatedAt,
	)
	return i, err
}
