package handlers

import (
	"database/sql"
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"
	"time"
)

type UserProfileHandler struct {
	queries *sqlc.Queries
}

func NewUserProfileHandler(q *sqlc.Queries) *UserProfileHandler {
	return &UserProfileHandler{
		queries: q,
	}
}

func (h *UserProfileHandler) HandleUserProfiles(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Ensure only /user-profiles endpoint is handled
	if len(parts) != 2 || parts[1] != "user-profiles" {
		http.Error(w, "Invalid URL - must be '/user-profiles'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllUserProfiles(w, r)
	case http.MethodPost:
		h.CreateUserProfile(w, r)
	default:
		http.Error(w, "Method not allowed - only GET and POST allowed at /user-profiles", http.StatusMethodNotAllowed)
	}
}

func (h *UserProfileHandler) GetAllUserProfiles(w http.ResponseWriter, r *http.Request) {}

func (h *UserProfileHandler) CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID       int32  `json:"user_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		DateOfBirth  string `json:"date_of_birth"`
		Gender       string `json:"gender"`
		HeightInches int32  `json:"height_inches"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dob, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	profile, err := h.queries.CreateUserProfile(r.Context(), sqlc.CreateUserProfileParams{
		UserID: sql.NullInt32{
			Int32: request.UserID,
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
		DateOfBirth: sql.NullTime{
			Time:  dob,
			Valid: true,
		},
		Gender: sql.NullString{
			String: request.Gender,
			Valid:  true,
		},
		HeightInches: sql.NullInt32{
			Int32: request.HeightInches,
			Valid: true,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
