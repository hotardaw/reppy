package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
)

type WorkoutSetByExerciseHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}

func NewWorkoutSetByExerciseHandler(q *sqlc.Queries, jwtSecret []byte) *WorkoutSetByExerciseHandler {
	return &WorkoutSetByExerciseHandler{
		queries:   q,
		jwtSecret: jwtSecret,
	}
}

func (h *WorkoutSetByExerciseHandler) HandleWorkoutSetsByExercise(w http.ResponseWriter, r *http.Request) {
	// get workout_id from URL: "/workouts/{workout_id}/exercises/{exercise_id}/sets"
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

	exerciseID, err := strconv.ParseInt(pathParts[4], 10, 32)
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.DeleteWorkoutSetsByExercise(w, r, workoutID, exerciseID) // "/workouts/{workout_id}/exercises/{exercise_id}/sets"
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// no set ID in path
func (h *WorkoutSetByExerciseHandler) DeleteWorkoutSetsByExercise(w http.ResponseWriter, r *http.Request, workoutID, exerciseID int64) {
	params := sqlc.DeleteWorkoutSetsByExerciseParams{
		WorkoutID:  int32(workoutID),
		ExerciseID: int32(exerciseID),
	}

	err := h.queries.DeleteWorkoutSetsByExercise(r.Context(), params)
	if err != nil {
		response.SendError(w, "Failed to delete workout sets by exercise", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, nil, http.StatusNoContent)
}
