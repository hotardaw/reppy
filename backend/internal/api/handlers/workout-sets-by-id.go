package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"strconv"
	"strings"
)

type WorkoutSetByIDHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}

func NewWorkoutSetByIDHandler(q *sqlc.Queries, jwtSecret []byte) *WorkoutSetByIDHandler {
	return &WorkoutSetByIDHandler{
		queries:   q,
		jwtSecret: jwtSecret,
	}
}

type UpdateWorkoutSetByIDRequest struct {
	Reps             *int32  `json:"reps"`
	ResistanceValue  *string `json:"resistance_value"`
	ResistanceType   *string `json:"resistance_type"`
	ResistanceDetail *string `json:"resistance_detail"`
	RPE              *string `json:"rpe"`
	Notes            *string `json:"notes"`
}

func (h *WorkoutSetByIDHandler) HandleWorkoutSets(w http.ResponseWriter, r *http.Request) {
	// "/workouts/{workout_id}/exercises/{exercise_id}/sets/{set_number}""
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 7 {
		response.SendError(w, "Invalid path URL", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(pathParts[2], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	exerciseID, err := strconv.ParseInt(pathParts[4], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	setNumber, err := strconv.ParseInt(pathParts[6], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid set number", http.StatusBadRequest)
		return
	}

	switch r.Method { // "/workouts/{workout_id}/exercises/{exercise_id}/sets/{set_number}""
	case http.MethodPatch:
		h.UpdateWorkoutSetByID(w, r, workoutID, exerciseID, setNumber)
	case http.MethodDelete:
		h.DeleteWorkoutSetByID(w, r, workoutID, exerciseID, setNumber)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// set ID in path - if something doesn't work make sure everything is SetByID-related
func (h *WorkoutSetByIDHandler) UpdateWorkoutSetByID(w http.ResponseWriter, r *http.Request, workoutID, exerciseID, setNumber int64) {
	var request UpdateWorkoutSetByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateWorkoutSetByIDParams{
		WorkoutID:        int32(workoutID),
		ExerciseID:       int32(exerciseID),
		SetNumber:        int32(setNumber),
		Reps:             utils.ToNullInt32(request.Reps),
		ResistanceValue:  utils.ToNullStringFromStringPtr(request.ResistanceValue),
		ResistanceType:   utils.ToNullResistanceTypeEnumFromStringPtr(request.ResistanceType),
		ResistanceDetail: utils.ToNullStringFromStringPtr(request.ResistanceDetail),
		Rpe:              utils.ToNullStringFromStringPtr(request.RPE),
		Notes:            utils.ToNullStringFromStringPtr(request.Notes),
	}

	set, err := h.queries.UpdateWorkoutSetByID(r.Context(), params)
	if err != nil {
		response.SendError(w, "Failed to update workout set", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, set, http.StatusCreated)
}

func (h *WorkoutSetByIDHandler) DeleteWorkoutSetByID(w http.ResponseWriter, r *http.Request, workoutID, exerciseID, setNumber int64) {
	params := sqlc.DeleteWorkoutSetByIDParams{
		WorkoutID:  int32(workoutID),
		ExerciseID: int32(exerciseID),
		SetNumber:  int32(setNumber),
	}

	err := h.queries.DeleteWorkoutSetByID(r.Context(), params)
	if err != nil {
		response.SendError(w, "Failed to delete workout set", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, nil, http.StatusNoContent)
}
