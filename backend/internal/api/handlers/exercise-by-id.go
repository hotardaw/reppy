package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
)

type ExerciseByIDHandler struct {
	queries *sqlc.Queries
}

func NewExerciseByIDHandler(q *sqlc.Queries) *ExerciseByIDHandler {
	return &ExerciseByIDHandler{queries: q}
}

func (h *ExerciseByIDHandler) HandleExercisesByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/exercises/")
	if idStr == "" {
		response.SendError(w, "Missing exercise ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetExerciseByID(w, r, int32(id))
	case http.MethodDelete:
		h.DeleteExercise(w, r, int32(id))
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/exercises/"
func (h *ExerciseByIDHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request, id int32) {
	exercise, err := h.queries.GetExerciseById(r.Context(), id)
	if err != nil {
		response.SendError(w, "Exercise not found", http.StatusNotFound)
		return
	}

	response.SendSuccess(w, exercise)
}

// "/exercises/"
func (h *ExerciseByIDHandler) DeleteExercise(w http.ResponseWriter, r *http.Request, id int32) {
	deletedExercise, err := h.queries.DeleteExercise(r.Context(), id)
	if err != nil {
		response.SendError(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message":  "Exercise deleted successfully",
		"exercise": deletedExercise,
	}, http.StatusOK) // Not StatusNoContent bc this is a soft delete
}
