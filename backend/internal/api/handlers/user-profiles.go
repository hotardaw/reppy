package handlers

import (
	"database/sql"
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
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

	switch r.Method {
	case http.MethodGet:
		if len(parts) > 2 {
			h.GetUserProfile(w, r, parts)
		} else {
			h.GetAllUserProfiles(w, r /*, parts*/)
		}
	case http.MethodPost:
		if len(parts) > 2 {
			http.Error(w, "Invalid URL for POST request", http.StatusBadRequest)
			return
		}
		h.CreateUserProfile(w, r)
	case http.MethodPut:
		if len(parts) > 2 {
			http.Error(w, "Invalid URL for PUT request", http.StatusBadRequest)
			return
		}
		h.UpdateUserProfile(w, r /*, parts*/)
	case http.MethodDelete:
		if len(parts) > 2 {
			http.Error(w, "Invalid URL for PUT request", http.StatusBadRequest)
			return
		}
		h.DeleteUserProfile(w, r /*, parts*/)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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

func (h *UserProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userProfile, err := h.queries.GetUserProfile(r.Context(), sql.NullInt32{Int32: int32(id), Valid: true})
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
func (h *UserProfileHandler) GetAllUserProfiles(w http.ResponseWriter, r *http.Request) {}
func (h *UserProfileHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request)  {}
func (h *UserProfileHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request)  {}
