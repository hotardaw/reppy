// GET (all & one by date) and POST only
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go-reppy/backend/internal/api/middleware"
	"go-reppy/backend/internal/api/response"
	"go-reppy/backend/internal/api/utils"
	"go-reppy/backend/internal/database/sqlc"
)

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

type CreateWorkoutRequest struct {
	Title       string           `json:"title"`             // name of workout, e.g. "Upper Body 2"
	WorkoutDate utils.CustomDate `json:"clientworkoutdate"` // YYYY-MM-DD
}

type WorkoutResponse struct {
	WorkoutID   int32          `json:"WorkoutID"`
	WorkoutDate string         `json:"WorkoutDate"` // post-db query, pre-response, sent as string
	Title       sql.NullString `json:"Title"`
	CreatedAt   sql.NullTime   `json:"CreatedAt"`
}

func (h *WorkoutHandler) HandleWorkouts(w http.ResponseWriter, r *http.Request) {
	dateQueryParams := r.URL.Query()
	if _, exists := dateQueryParams["date"]; exists {
		if r.Method == http.MethodGet {
			h.GetWorkoutByUserIDAndDate(w, r)
			return
		}
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.CreateWorkout(w, r)
	case http.MethodGet:
		h.GetAllWorkoutsForUser(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/workouts"
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

	// must convert to user's tz before returning data to user
	var resp []WorkoutResponse
	for _, workout := range workouts {
		// convert UTC to client tz
		clientTime, err := utils.FromUTCToClientTimezone(workout.WorkoutDate, r)
		if err != nil {
			response.SendError(w, "Timezone conversion error", http.StatusBadRequest)
			return
		}

		resp = append(resp, WorkoutResponse{
			WorkoutID:   workout.WorkoutID,
			WorkoutDate: clientTime.Format("2006-01-02"),
			Title:       workout.Title,
			CreatedAt:   workout.CreatedAt,
		})
	}

	response.SendSuccess(w, resp)
}

// "/workouts?date=2024-01-05"
func (h *WorkoutHandler) GetWorkoutByUserIDAndDate(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		response.SendError(w, "Date parameter is required in format: 2024-12-25", http.StatusBadRequest)
		return
	}

	workoutDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, "Invalid date format - use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	utcTime, err := utils.FromClientTimezoneToUTC(workoutDate, r)
	if err != nil {
		response.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout, err := h.queries.GetWorkoutByUserIDAndDate(r.Context(), sqlc.GetWorkoutByUserIDAndDateParams{
		UserID:      utils.ToNullInt32(userID),
		WorkoutDate: utcTime,
	})
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

// "/workouts"
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var request CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutDate := request.WorkoutDate.ToTime()

	utcTime, err := utils.FromClientTimezoneToUTC(workoutDate, r)
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
