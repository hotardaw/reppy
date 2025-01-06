// For /user-profiles/{user_id} endpoint
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
)

type UserProfileByIDHandler struct {
	queries *sqlc.Queries
}

func NewUserProfileByIDHandler(q *sqlc.Queries) *UserProfileByIDHandler {
	return &UserProfileByIDHandler{
		queries: q,
	}
}

func (h *UserProfileByIDHandler) HandleUserProfilesByID(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Ensure only /user-profiles/{id} endpoint is handled
	if len(parts) != 3 {
		response.SendError(w, "Invalid URL format - must be '/user-profiles/{user_id}'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserProfile(w, r, parts)
	case http.MethodPatch:
		h.UpdateUserProfile(w, r, parts)
	case http.MethodDelete:
		h.DeleteUserProfile(w, r, parts)
	default:
		response.SendError(w, "Method not allowed - only GET, PUT, and DELETE", http.StatusMethodNotAllowed)
		return
	}
}

func (h *UserProfileByIDHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userProfile, err := h.queries.GetUserProfile(r.Context(), utils.ToNullInt32(id))
	if err == sql.ErrNoRows {
		response.SendError(w, "User profile not found", http.StatusNotFound)
		return
	} else if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfile)
}

func (h *UserProfileByIDHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		FirstName    string `json:"first_name,omitempty"`
		LastName     string `json:"last_name,omitempty"`
		HeightInches int32  `json:"height_inches,omitempty"`
		WeightPounds int32  `json:"weight_pounds,omitempty"`
		Gender       string `json:"gender,omitempty"`
		DateOfBirth  string `json:"date_of_birth,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateUserProfileParams := sqlc.UpdateUserProfileParams{
		UserID:       utils.ToNullInt32(id),
		FirstName:    utils.ToNullString(request.FirstName),
		LastName:     utils.ToNullString(request.LastName),
		HeightInches: utils.ToNullInt32(request.HeightInches),
		WeightPounds: utils.ToNullInt32(request.WeightPounds),
		Gender:       utils.ToNullString(request.Gender),
	}

	// If provided, handle DOB
	if request.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", request.DateOfBirth)
		if err != nil {
			response.SendError(w, "Invalid date format for date_of_birth; use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		updateUserProfileParams.DateOfBirth = utils.ToNullTime(dob)
	}

	userProfile, err := h.queries.UpdateUserProfile(r.Context(), updateUserProfileParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfile)
}

func (h *UserProfileByIDHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userProfile, err := h.queries.DeleteUserProfile(r.Context(), utils.ToNullInt32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User profile not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfile, http.StatusOK) // Not StatusNoContent bc this is a soft delete
}
