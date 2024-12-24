// GET (ID only), PATCH, DELETE
package handlers

import (
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"
)

// GET, PATCH, DELETE w/ ID
type WorkoutByIDHandler struct {
	queries *sqlc.Queries
}

func NewWorkoutByIDHandler(q *sqlc.Queries) *WorkoutByIDHandler {
	return &WorkoutByIDHandler{
		queries: q,
	}
}

func (h *WorkoutByIDHandler) HandleWorkoutsByID(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	if len(parts) != 3 {
		response.SendError(w, "Invalid URL - must be '/workouts/{workout_id}'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetWorkoutByID(w, r)
	case http.MethodPatch:
		h.UpdateWorkoutByID(w, r)
	case http.MethodDelete:
		h.DeleteWorkoutByID(w, r)
	}
}

func (h *WorkoutByIDHandler) GetWorkoutByID(w http.ResponseWriter, r *http.Request)    {}
func (h *WorkoutByIDHandler) UpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {}
func (h *WorkoutByIDHandler) DeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {}
