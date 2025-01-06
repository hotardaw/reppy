package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
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
	WorkoutID    int32 `json:"workout_id"`
	ExerciseID   int32 `json:"exercise_id"`
	NumberOfSets int32 `json:"number_of_sets"` // request JSON data will be duplicated this # of times for each set
	// optional fields (nil if absent):
	Reps             *int32  `json:"reps"`
	ResistanceValue  *string `json:"resistance_value"`
	ResistanceType   *string `json:"resistance_type"`
	ResistanceDetail *string `json:"resistance_detail"`
	RPE              *string `json:"rpe"`
	Notes            *string `json:"notes"`
}

/*
sample req body to test auto-incr in postman:
{
  "workout_id": 1,
  "exercise_id": 1,
  "number_of_sets": 3,
  "reps": 10,
  "resistance_value": "135.5",
  "resistance_type": "weight",
  "resistance_detail": "barbell",
  "rpe": "8"
}

sample minimal request:
{
    "workout_id": 1,
    "exercise_id": 1,
    "number_of_sets": 3
}
*/

func (h *WorkoutSetHandler) HandleWorkoutSets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateWorkoutSets(w, r)
	}
}

func (h *WorkoutSetHandler) CreateWorkoutSets(w http.ResponseWriter, r *http.Request) {
	var request CreateWorkoutSetsRequest
	fmt.Println("Request: ", request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.NumberOfSets <= 0 {
		response.SendError(w, "Number of sets must be greater than 0", http.StatusBadRequest)
		return
	}

	params := sqlc.CreateWorkoutSetsParams{
		Column1: request.WorkoutID,                    // workout_id
		Column2: request.ExerciseID,                   // exercise_id
		Column3: make([]int32, request.NumberOfSets),  // set_number
		Column4: make([]int32, request.NumberOfSets),  // reps
		Column5: make([]string, request.NumberOfSets), // resistance_value
		Column6: make([]string, request.NumberOfSets), // resistance_type
		Column7: make([]string, request.NumberOfSets), // resistance_detail
		Column8: make([]string, request.NumberOfSets), // rpe
		Column9: make([]string, request.NumberOfSets), // notes
	}

	for i := int32(0); i < request.NumberOfSets; i++ {
		params.Column3[i] = i + 1 // auto-incr set nums starting with 1

		if request.Reps != nil {
			params.Column4[i] = *request.Reps
		}
		if request.ResistanceValue != nil {
			params.Column5[i] = *request.ResistanceValue
		}
		if request.ResistanceType != nil {
			params.Column6[i] = *request.ResistanceType
		}
		if request.ResistanceDetail != nil {
			params.Column7[i] = *request.ResistanceDetail
		}
		if request.RPE != nil {
			params.Column8[i] = *request.RPE
		}
		if request.Notes != nil {
			params.Column9[i] = *request.Notes
		}
	}

	fmt.Println("Params: ", params)

	sets, err := h.queries.CreateWorkoutSets(r.Context(), params)
	if err != nil {
		response.SendError(w, "Failed to create workout set(s)"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Sets: ", sets)

	response.SendSuccess(w, sets, http.StatusCreated)
}
