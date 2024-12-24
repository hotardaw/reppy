package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type ExerciseHandler struct {
	queries *sqlc.Queries
}

func NewExerciseHandler(q *sqlc.Queries) *ExerciseHandler {
	return &ExerciseHandler{
		queries: q,
	}
}

func (h *ExerciseHandler) HandleExercises(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	switch r.Method {
	case http.MethodGet: // Handle both GET all '/exercises' and GET '/exercises/{id}'
		if len(parts) == 2 {
			h.GetAllExercises(w, r)
		} else if len(parts) == 3 {
			h.GetExerciseByID(w, r)
		} else {
			response.SendError(w, "Invalid URL format for GET request", http.StatusBadRequest)
		}
	case http.MethodPost: // Handle '/exercises'
		h.CreateExercise(w, r)
	case http.MethodDelete: // Handle '/exercises/{id}'
		if len(parts) != 3 {
			response.SendError(w, "DELETE requests must be to '/exercises/{id}'", http.StatusBadRequest)
			return
		}
		h.DeleteExercise(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(path.Clean(r.URL.Path), "/")

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.queries.GetExerciseById(r.Context(), int32(id))
	if err != nil {
		response.SendError(w, "Muscle not found", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercise)
}

func (h *ExerciseHandler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {
	exerciseName := r.URL.Query().Get("name") // '/exercises?name=Bench%20Press'
	if exerciseName == "" {
		response.SendError(w, "Exercise name is required for GET requests", http.StatusBadRequest)
		return
	}

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
	var request struct {
		ExerciseName string `json:"exercise_name"`
		Description  string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.queries.CreateExercise(r.Context(), sqlc.CreateExerciseParams{
		ExerciseName: request.ExerciseName,
		Description:  utils.ToNullString(request.Description),
	})
	if err != nil {
		response.SendError(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercise, http.StatusCreated)
}

func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	if len(parts) != 3 { // /exercises/{id} should have 3 parts
		response.SendError(w, "Invalid URL format - must be '/exercises/{exercise_id}'", http.StatusBadRequest)
		return
	}

	exerciseId, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.queries.DeleteExercise(r.Context(), int32(exerciseId))
	if err != nil {
		response.SendError(w, "Failed to delete muscle", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, exercise, http.StatusNoContent)
}
