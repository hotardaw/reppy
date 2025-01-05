package handlers

import (
	"encoding/json"
	"net/http"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
)

type ExerciseHandler struct {
	queries *sqlc.Queries
}

type CreateExerciseRequest struct {
	ExerciseName string `json:"exercise_name"`
	Description  string `json:"description"`
}

func NewExerciseHandler(q *sqlc.Queries) *ExerciseHandler {
	return &ExerciseHandler{queries: q}
}

func (h *ExerciseHandler) HandleExercises(w http.ResponseWriter, r *http.Request) {
	if name := r.URL.Query().Get("name"); name != "" {
		h.GetExerciseByName(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllExercises(w, r)
	case http.MethodPost:
		h.CreateExercise(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExerciseHandler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {
	exerciseName := r.URL.Query().Get("name") // '/exercises?name=Bench%20Press'

	exercise, err := h.queries.GetExerciseByName(r.Context(), exerciseName)
	if err != nil {
		response.SendError(w, "Exercise not found", http.StatusInternalServerError)
	}

	response.SendSuccess(w, exercise)
}

func (h *ExerciseHandler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.queries.GetAllExercises(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercises)
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var request CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createExerciseParams := sqlc.CreateExerciseParams{
		ExerciseName: request.ExerciseName,
		Description:  utils.ToNullString(request.Description),
	}
	exercise, err := h.queries.CreateExercise(r.Context(), createExerciseParams)
	if err != nil {
		response.SendError(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercise, http.StatusCreated)
}
