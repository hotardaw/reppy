package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
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
	ExerciseID       int32   `json:"exercise_id"`
	Reps             *int32  `json:"reps"`
	ResistanceValue  *string `json:"resistance_value"`
	ResistanceType   *string `json:"resistance_type"`
	ResistanceDetail *string `json:"resistance_detail"`
	RPE              *string `json:"rpe"`
	Notes            *string `json:"notes"`
}

func (h *WorkoutSetByIDHandler) HandleWorkoutSetByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.SendError(w, "Invalid path URL", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(pathParts[2], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	overallSetNumber, err := strconv.ParseInt(pathParts[4], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid set number", http.StatusBadRequest)
		return
	}

	switch r.Method { // "/workouts/3/workout-sets/7"
	case http.MethodPatch:
		h.UpdateWorkoutSetByID(w, r, int32(workoutID), int32(overallSetNumber))
	case http.MethodDelete:
		h.DeleteWorkoutSetByID(w, r, int32(workoutID), int32(overallSetNumber))
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/workouts/3/workout-sets/7"
func (h *WorkoutSetByIDHandler) UpdateWorkoutSetByID(w http.ResponseWriter, r *http.Request, workoutID, overallSetNumber int32) {
	var request UpdateWorkoutSetByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	set, err := h.queries.UpdateWorkoutSetByID(r.Context(), sqlc.UpdateWorkoutSetByIDParams{
		WorkoutID:               workoutID,
		OverallWorkoutSetNumber: overallSetNumber,
		Reps:                    utils.ToNullInt32FromIntPtr(request.Reps),
		ResistanceValue:         utils.ToNullStringFromStringPtr(request.ResistanceValue),
		ResistanceType:          utils.ToNullResistanceTypeEnumFromStringPtr(request.ResistanceType),
		ResistanceDetail:        utils.ToNullStringFromStringPtr(request.ResistanceDetail),
		Rpe:                     utils.ToNullStringFromStringPtr(request.RPE),
		Notes:                   utils.ToNullStringFromStringPtr(request.Notes),
	})
	if err != nil {
		response.SendError(w, "Failed to update workout set", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, set)
}

// "/workouts/3/workout-sets/7"
func (h *WorkoutSetByIDHandler) DeleteWorkoutSetByID(w http.ResponseWriter, r *http.Request, workoutID, overallSetNumber int32) {
	err := h.queries.DeleteWorkoutSetByID(r.Context(), sqlc.DeleteWorkoutSetByIDParams{
		WorkoutID:               workoutID,
		OverallWorkoutSetNumber: overallSetNumber,
	})
	if err != nil {
		response.SendError(w, "Failed to delete workout set", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, nil, http.StatusOK)
}
