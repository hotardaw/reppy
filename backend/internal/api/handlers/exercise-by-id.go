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
	id := strings.TrimPrefix(r.URL.Path, "/exercises/")
	if id == "" {
		response.SendError(w, "Missing exercise ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetExerciseByID(w, r, id)
	case http.MethodDelete:
		h.DeleteExercise(w, r, id)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExerciseByIDHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.queries.GetExerciseById(r.Context(), int32(id))
	if err != nil {
		response.SendError(w, "Exercise not found", http.StatusNotFound)
		return
	}

	response.SendSuccess(w, exercise)
}

func (h *ExerciseByIDHandler) DeleteExercise(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.queries.DeleteExercise(r.Context(), int32(id))
	if err != nil {
		response.SendError(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercise, http.StatusNoContent)
}
