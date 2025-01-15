package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
  "exercise_id": 1,
  "number_of_sets": 3
}
*/

// "/workouts/{workout_id}/workout-sets"
func (h *WorkoutSetHandler) HandleWorkoutSets(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.SendError(w, "Invalid path URL", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(pathParts[2], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	switch r.Method { // "/workouts/{workout_id}/sets"
	case http.MethodPost:
		h.CreateWorkoutSets(w, r, int32(workoutID))
	case http.MethodGet:
		h.GetAllWorkoutSets(w, r, int32(workoutID))
	case http.MethodDelete:
		h.DeleteAllWorkoutSets(w, r, int32(workoutID))
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/workouts/3/workout-sets"
func (h *WorkoutSetHandler) CreateWorkoutSets(w http.ResponseWriter, r *http.Request, workoutID int32) {
	var request CreateWorkoutSetsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.NumberOfSets <= 0 {
		response.SendError(w, "Number of sets must be greater than 0", http.StatusBadRequest)
		return
	}

	if request.ExerciseID <= 0 {
		response.SendError(w, "Exercise ID must be provided", http.StatusBadRequest)
		return
	}

	params := sqlc.CreateWorkoutSetsParams{
		Column1: workoutID,                            // workout_id
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

	sets, err := h.queries.CreateWorkoutSets(r.Context(), params)
	if err != nil {
		response.SendError(w, "Failed to create workout set(s)"+err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, sets, http.StatusCreated)
}

// "/workouts/{workout_id}/sets"
func (h *WorkoutSetHandler) GetAllWorkoutSets(w http.ResponseWriter, r *http.Request, workoutID int32) {
	allWorkoutSets, err := h.queries.GetAllWorkoutSets(r.Context(), workoutID)
	if err != nil {
		response.SendError(w, "All workout sets not found", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, allWorkoutSets)
}

// "/workouts/{workout_id}/sets"
func (h *WorkoutSetHandler) DeleteAllWorkoutSets(w http.ResponseWriter, r *http.Request, workoutID int32) {
	err := h.queries.DeleteAllWorkoutSets(r.Context(), workoutID)
	if err != nil {
		response.SendError(w, "Failed to delete all workout sets", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, nil, http.StatusNoContent)
}
