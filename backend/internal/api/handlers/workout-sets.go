package handlers

import (
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
)

type WorkoutSetHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}

func NewWorkoutSetHandler(q *sqlc.Queries, jwtSecret []byte) *WorkoutSetHandler {
	return &WorkoutSetHandler{
		queries:   q,
		jwtSecret: jwtSecret,
	}
}

// req, resp structs
type CreateWorkoutSetsRequest struct {
	Sets []struct {
		WorkoutID        int32    `json:"workout_id"`
		ExerciseID       int32    `json:"exercise_id"`
		SetNumber        int32    `json:"set_number"`
		Reps             *int32   `json:"reps,omitempty"`
		ResistanceType   *string  `json:"resistance_type,omitempty"`
		ResistanceValue  *float32 `json:"resistance_value,omitempty"`
		ResistanceDetail *string  `json:"resistance_detail,omitempty"`
		RPE              *float32 `json:"rpe,omitempty"`
		Notes            *string  `json:"notes,omitempty"`
	} `json:"sets"`
}

func (h *WorkoutSetHandler) HandleWorkoutSets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateWorkoutSets(w, r)
	}
}

// Used for batch inserts and single inserts.
func (h *WorkoutSetHandler) CreateWorkoutSets(w http.ResponseWriter, r *http.Request) {

}

// Helper for converting request JSON to SQLc params
func (r *CreateWorkoutSetsRequest) ToSQLCParams() sqlc.CreateWorkoutSetsParams {
	n := len(r.Sets)
	params := sqlc.CreateWorkoutSetsParams{
		Column1: make([]int32, n),  // workout_ids
		Column2: make([]int32, n),  // exercise_ids
		Column3: make([]int32, n),  // set_numbers
		Column4: make([]int32, n),  // reps
		Column5: make([]string, n), // resistance_values
		Column6: make([]string, n), // resistance_types
		Column7: make([]string, n), // resistance_details
		Column8: make([]string, n), // rpes
		Column9: make([]string, n), // notes
	}

	for i, set := range r.Sets {
		params.Column1[i] = set.WorkoutID
		params.Column2[i] = set.ExerciseID
		params.Column3[i] = set.SetNumber
		// the ".Int32" and ".String" bits select the int32 & string values from the returned struct, implicitly discarding the "Valid" field.
		params.Column4[i] = utils.NullIntFromIntPtr(set.Reps).Int32
		params.Column5[i] = utils.NullStringFromFloat32Ptr(set.ResistanceValue).String
		params.Column6[i] = utils.NullStringFromStringPtr(set.ResistanceType).String
		params.Column7[i] = utils.NullStringFromStringPtr(set.ResistanceDetail).String
		params.Column8[i] = utils.NullStringFromFloat32Ptr(set.RPE).String
		params.Column9[i] = utils.NullStringFromStringPtr(set.Notes).String
	}

	return params
}
