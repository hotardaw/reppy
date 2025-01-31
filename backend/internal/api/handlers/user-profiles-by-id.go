// For /user-profiles/{user_id} endpoint
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-fitstat/backend/internal/api/response"
	"go-fitstat/backend/internal/api/utils"
	"go-fitstat/backend/internal/database/sqlc"
)

type UserProfileByIDHandler struct {
	queries *sqlc.Queries
}

func NewUserProfileByIDHandler(q *sqlc.Queries) *UserProfileByIDHandler {
	return &UserProfileByIDHandler{
		queries: q,
	}
}

type UpdateUserProfileRequest struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	HeightInches int32  `json:"height_inches,omitempty"`
	WeightPounds int32  `json:"weight_pounds,omitempty"`
	Gender       string `json:"gender,omitempty"`
	DateOfBirth  string `json:"date_of_birth,omitempty"`
}

func (h *UserProfileByIDHandler) HandleUserProfilesByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid user ID in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserProfile(w, r, id)
	case http.MethodPatch:
		h.UpdateUserProfile(w, r, id)
	case http.MethodDelete:
		h.DeleteUserProfile(w, r, id)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/user-profiles/{id}"
func (h *UserProfileByIDHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, id int) {
	userProfile, err := h.queries.GetUserProfile(r.Context(), utils.ToNullInt32(id))
	if err == sql.ErrNoRows {
		response.SendError(w, "User profile not found", http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendError(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfile)
}

// "/user-profiles/{id}"
func (h *UserProfileByIDHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request, id int) {
	var request UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateUserProfileParams{
		UserID:       utils.ToNullInt32(id),
		FirstName:    utils.ToNullString(request.FirstName),
		LastName:     utils.ToNullString(request.LastName),
		HeightInches: utils.ToNullInt32(request.HeightInches),
		WeightPounds: utils.ToNullInt32(request.WeightPounds),
		Gender:       utils.ToNullString(request.Gender),
	}

	if request.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", request.DateOfBirth)
		if err != nil {
			response.SendError(w, "Invalid date format for date_of_birth (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		params.DateOfBirth = utils.ToNullTime(dob)
	}

	userProfile, err := h.queries.UpdateUserProfile(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		response.SendError(w, "User profile not found", http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendError(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfile)
}

// "/user-profiles/{id}"
func (h *UserProfileByIDHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request, id int) {
	deletedUserProfile, err := h.queries.DeleteUserProfile(r.Context(), utils.ToNullInt32(id))
	if errors.Is(err, sql.ErrNoRows) {
		response.SendError(w, "User profile not found", http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendError(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message": "User profile (soft-) deleted successfully",
		"id":      deletedUserProfile,
	}, http.StatusOK) // Not StatusNoContent bc this is a soft delete)
}
