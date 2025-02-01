package seeds

import (
	"context"
	"fmt"
	"strings"

	"go-reppy/backend/internal/api/utils"
	"go-reppy/backend/internal/database/sqlc"
)

type TestWorkoutSets struct { // Fields with pointers are nil-able
	WorkoutID        int32
	ExerciseID       int32
	SetNumber        int32 // Incrementing logic will need to be handled in app layer
	Reps             *int32
	ResistanceType   *string  // 'weight', 'band', 'bodyweight'
	ResistanceValue  *float32 // weight in pounds/kgs
	ResistanceDetail *string  // "blue band", "wide grip", "olympic bar", etc.
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
			Reps:             utils.IntPtr(8),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Easy af"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        2,
			Reps:             utils.IntPtr(8),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Less easy this time"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        3,
			Reps:             utils.IntPtr(7),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Barely 7"),
		},
	}
}

// Failed to seed test data: failed to seed workouts: failed to seed workout-sets 1, 4, 1: pq: invalid input value for enum resistance_type_enum: "Bodyweight, +0 lbs"
func SeedWorkoutSets(queries *sqlc.Queries) error {
	sets := GetTestWorkoutSets()

	// group sets by workout_id & exercise_id
	groupedSets := make(map[struct{ workoutID, exerciseID int32 }][]TestWorkoutSets)
	for _, set := range sets {
		key := struct{ workoutID, exerciseID int32 }{set.WorkoutID, set.ExerciseID}
		groupedSets[key] = append(groupedSets[key], set)
	}

	// insert groups individually
	for key, sets := range groupedSets {
		setNumbers := make([]int32, len(sets))
		reps := make([]int32, len(sets))
		resistanceValues := make([]string, len(sets))
		resistanceTypes := make([]string, len(sets)) // changed from []sqlc.ResistanceTypeEnum
		resistanceDetails := make([]string, len(sets))
		rpes := make([]string, len(sets))
		notes := make([]string, len(sets))

		for i, set := range sets {
			setNumbers[i] = set.SetNumber
			if set.Reps != nil {
				reps[i] = *set.Reps
			}
			if set.ResistanceValue != nil {
				resistanceValues[i] = fmt.Sprintf("%.1f", *set.ResistanceValue)
			}
			if set.ResistanceType != nil {
				resistanceTypes[i] = strings.ToLower(*set.ResistanceType)
			}
			if set.ResistanceDetail != nil {
				resistanceDetails[i] = *set.ResistanceDetail
			}
			if set.RPE != nil {
				rpes[i] = fmt.Sprintf("%.1f", *set.RPE)
			}
			if set.Notes != nil {
				notes[i] = *set.Notes
			}
		}

		_, err := queries.CreateWorkoutSets(context.Background(), sqlc.CreateWorkoutSetsParams{
			Column1: key.workoutID,  // single value
			Column2: key.exerciseID, // single value
			Column3: setNumbers,
			Column4: reps,
			Column5: resistanceValues,
			Column6: resistanceTypes, // now just []string, SQL handles conversion
			Column7: resistanceDetails,
			Column8: rpes,
			Column9: notes,
		})
		if err != nil {
			return fmt.Errorf("failed to seed workout-sets for workout %d, exercise %d: %v",
				key.workoutID, key.exerciseID, err)
		}
	}

	fmt.Println("Successfully seeded WORKOUT-SETS table")
	return nil
}
