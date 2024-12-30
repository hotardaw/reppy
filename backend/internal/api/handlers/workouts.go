// GET (all & one by date) and POST only
package handlers

import (
	"context"
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

// GET (all, one by date), POST
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

func (h *WorkoutHandler) HandleWorkouts(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	if len(parts) == 2 { // "/workouts"
		switch r.Method {
		case http.MethodGet:
			h.GetAllWorkoutsForUser(w, r)
		case http.MethodPost:
			h.CreateWorkout(w, r)
		default:
			response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 3 && parts[2] == "date" { // "workouts/date"
		switch r.Method {
		case http.MethodGet:
			h.GetWorkoutByUserIDAndDate(w, r)
		default:
			response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	response.SendError(w, "Invalid URL", http.StatusBadRequest)
}

func (h *WorkoutHandler) GetWorkoutByUserIDAndDate(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		WorkoutDate time.Time `json:"clientworkoutdate"`
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

	params := sqlc.GetWorkoutByUserIDAndDateParams{
		UserID:      utils.ToNullInt32(userID),
		WorkoutDate: utcTime,
	}

	workout, err := h.queries.GetWorkoutByUserIDAndDate(r.Context(), params)
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

func (h *WorkoutHandler) GetAllWorkoutsForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
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
		Title       string    `json:"title"` // name of workout, e.g. "Upper Body 2"
		WorkoutDate time.Time `json:"clientworkoutdate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	utcTime, err := utils.FromClientTimezoneToUTC(request.WorkoutDate, r)
	if err != nil {
		response.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), sqlc.CreateWorkoutParams{
		UserID:      utils.ToNullInt32(userID),
		WorkoutDate: utcTime,
		Title:       utils.ToNullString(request.Title),
	})
	if err != nil {
		response.SendError(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout, http.StatusCreated)
}
