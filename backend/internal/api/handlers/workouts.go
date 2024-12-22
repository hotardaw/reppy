// GET (all) and POST only
package handlers

import (
	"database/sql"
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"
)

// GET, PATCH, DELETE w/ ID
type WorkoutHandler struct {
	queries *sqlc.Queries
}

func NewWorkoutHandler(q *sqlc.Queries) *WorkoutHandler {
	return &WorkoutHandler{
		queries: q,
	}
}

func (h *WorkoutHandler) HandleWorkoutsByID(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	if len(parts) != 2 {
		http.Error(w, "Invalid URL - must be '/workouts'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetWorkouts(w, r)
	case http.MethodPost:
		h.CreateWorkout(w, r)
	}
}

func (h *WorkoutHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {}

func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title string `json:"title"` // name of workout, e.g. "Upper Body 2"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), sqlc.CreateWorkoutParams{
		Title: sql.NullString{
			String: request.Title,
			Valid:  true,
		},
	})
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}
