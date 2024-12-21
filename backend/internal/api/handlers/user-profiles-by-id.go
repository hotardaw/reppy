// For /user-profiles/{user_id} endpoint
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
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
		http.Error(w, "Invalid URL format - must be '/user-profiles/{user_id}'", http.StatusBadRequest)
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
		http.Error(w, "Method not allowed - only GET, PUT, and DELETE", http.StatusMethodNotAllowed)
	}
}

func (h *UserProfileByIDHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userProfile, err := h.queries.GetUserProfile(r.Context(), sql.NullInt32{
		Int32: int32(id),
		Valid: true,
	})
	if err == sql.ErrNoRows {
		http.Error(w, "User profile not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserProfileByIDHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateUserProfileParams{
		UserID: sql.NullInt32{
			Int32: int32(id),
			Valid: true,
		},
		FirstName: sql.NullString{
			String: request.FirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: request.LastName,
			Valid:  true,
		},
		HeightInches: sql.NullInt32{
			Int32: request.HeightInches,
			Valid: true,
		},
		WeightPounds: sql.NullInt32{
			Int32: request.WeightPounds,
			Valid: true,
		},
		Gender: sql.NullString{
			String: request.Gender,
			Valid:  true,
		},
	}

	// Handle DOB, if provided
	if request.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", request.DateOfBirth)
		if err != nil {
			http.Error(w, "Invalid date format for date_of_birth; use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		params.DateOfBirth = sql.NullTime{
			Time:  dob,
			Valid: true,
		}
	}

	userProfile, err := h.queries.UpdateUserProfile(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserProfileByIDHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userProfile, err := h.queries.DeleteUserProfile(r.Context(), sql.NullInt32{
		Int32: int32(id),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}
