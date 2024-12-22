package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
)

type TestWorkoutSets struct { // Fields with pointers are nil-able
	WorkoutID        int32
	ExerciseID       int32
	SetNumber        int32 // Incrementing logic will need to be handled in app layer
	Reps             *int32
	ResistanceType   *string // 'weight', 'band', 'bodyweight'
	ResistanceValue  *int32  // weight in pounds/kgs
	ResistanceDetail *string // "blue band", "wide grip", "olympic bar", etc.
	RPE              *float32
	Notes            *string
}

func GetTestWorkoutSets() []TestWorkoutSets {
	return []TestWorkoutSets{
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  1,
		},
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  2,
		},
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  3,
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        1,
			Reps:             intPtr(8),
			ResistanceType:   strPtr("Bodyweight, +0 lbs"),
			ResistanceValue:  intPtr(0),
			ResistanceDetail: strPtr("Wide grip"),
			RPE:              float32Ptr(8.5),
			Notes:            strPtr("Easy af"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        2,
			Reps:             intPtr(8),
			ResistanceType:   strPtr("Bodyweight, +0 lbs"),
			ResistanceValue:  intPtr(0),
			ResistanceDetail: strPtr("Wide grip"),
			RPE:              float32Ptr(8.5),
			Notes:            strPtr("Less easy this time"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        1,
			Reps:             intPtr(7),
			ResistanceType:   strPtr("Bodyweight, +0 lbs"),
			ResistanceValue:  intPtr(0),
			ResistanceDetail: strPtr("Wide grip"),
			RPE:              float32Ptr(8.5),
			Notes:            strPtr("Barely 7"),
		},
	}
}

func SeedWorkoutExercises(queries *sqlc.Queries) error {
	for _, set := range GetTestWorkoutSets() {
		_, err := queries.CreateWorkoutSet(context.Background(), sqlc.CreateWorkoutSetParams{
			WorkoutID:        set.WorkoutID,
			ExerciseID:       set.ExerciseID,
			SetNumber:        set.SetNumber,
			Reps:             nullIntFromIntPtr(set.Reps),
			ResistanceValue:  nullIntFromIntPtr(set.ResistanceValue),
			ResistanceType:   nullResistanceTypeFromPtr(set.ResistanceType),
			ResistanceDetail: nullStringFromStringPtr(set.ResistanceDetail),
			Rpe:              nullStringFromFloat32Ptr(set.RPE),
			Notes:            nullStringFromStringPtr(set.Notes),
		})
		if err != nil {
			return fmt.Errorf("failed to seed workout-exercises %d, %d, %d: %v", set.WorkoutID, set.ExerciseID, set.SetNumber, err)
		}
	}
	return nil
}

func intPtr(i int32) *int32         { return &i }
func strPtr(s string) *string       { return &s }
func float32Ptr(f float32) *float32 { return &f }

// For converting pointers to SQL null types
func nullIntFromIntPtr(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func nullStringFromStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullStringFromFloat32Ptr(f *float32) sql.NullString {
	if f == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: fmt.Sprintf("%.1f", *f), Valid: true}
}

func nullResistanceTypeFromPtr(s *string) sqlc.NullResistanceTypeEnum {
	if s == nil {
		return sqlc.NullResistanceTypeEnum{}
	}
	return sqlc.NullResistanceTypeEnum{ResistanceTypeEnum: sqlc.ResistanceTypeEnum(*s), Valid: true}
}
