// GET (ID only), PATCH, DELETE
package handlers

import (
	"database/sql"
	"encoding/json"
	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"
	"time"
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
		h.GetWorkoutByIDForUser(w, r)
	case http.MethodPatch:
		h.UpdateWorkoutByID(w, r)
	case http.MethodDelete:
		h.DeleteWorkoutByID(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *WorkoutByIDHandler) GetWorkoutByIDForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	params := sqlc.GetWorkoutByIDForUserParams{
		WorkoutID: workoutID,
		UserID:    utils.ToNullInt32(userID),
	}

	workout, err := h.queries.GetWorkoutByIDForUser(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout)
}

func (h *WorkoutByIDHandler) UpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	var request struct {
		WorkoutDate time.Time `json:"clientworkoutdate"`
		Title       string    `json:"workouttitle"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	utcTime, err := utils.FromClientTimezoneToUTC(request.WorkoutDate, r)
	if err != nil {
		response.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateWorkoutParams{
		WorkoutDate: utcTime,
		Title:       utils.ToNullString(request.Title),
		WorkoutID:   workoutID,
		UserID:      utils.ToNullInt32(userID),
	}

	workout, err := h.queries.UpdateWorkout(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout)
}

func (h *WorkoutByIDHandler) DeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	params := sqlc.DeleteWorkoutParams{
		UserID:    utils.ToNullInt32(userID),
		WorkoutID: workoutID,
	}

	deletedWorkoutID, err := h.queries.DeleteWorkout(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message": "Workout deleted successfully",
		"id":      int(deletedWorkoutID),
	})
}
