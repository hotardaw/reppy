// GET (all) and POST only
package handlers

import (
	"context"
	"encoding/json"
	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"
)

// GET, PATCH, DELETE w/ ID
type WorkoutHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}

func NewWorkoutHandler(q *sqlc.Queries, jwtSecret []byte) *WorkoutHandler {
	return &WorkoutHandler{
		queries:   q,
		jwtSecret: jwtSecret,
	}
}

func (h *WorkoutHandler) HandleWorkoutsByID(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	if len(parts) != 2 {
		response.SendError(w, "Invalid URL - must be '/workouts'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllWorkoutsForUser(w, r)
	case http.MethodPost:
		h.CreateWorkout(w, r)
	}
}

func (h *WorkoutHandler) GetAllWorkoutsForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
	}

	workouts, err := h.queries.GetAllWorkoutsForUser(context.Background(), utils.ToNullInt32(userID))
	if err != nil {
		response.SendError(w, "Failed to get all user workouts", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workouts)
}

func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title string `json:"title"` // name of workout, e.g. "Upper Body 2"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), sqlc.CreateWorkoutParams{
		Title: utils.ToNullString(request.Title),
	})
	if err != nil {
		response.SendError(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout, http.StatusCreated)
}
