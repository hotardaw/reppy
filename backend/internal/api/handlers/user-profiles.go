package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
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
		response.SendError(w, "Invalid URL - must be '/user-profiles'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllUserProfiles(w, r)
	case http.MethodPost:
		h.CreateUserProfile(w, r)
	default:
		response.SendError(w, "Method not allowed - only GET and POST allowed at /user-profiles", http.StatusMethodNotAllowed)
	}
}

func (h *UserProfileHandler) GetAllUserProfiles(w http.ResponseWriter, r *http.Request) {
	userProfiles, err := h.queries.GetAllUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfiles)
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
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dob, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		response.SendError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	profile, err := h.queries.CreateUserProfile(r.Context(), sqlc.CreateUserProfileParams{
		UserID:       utils.ToNullInt32(request.UserID),
		FirstName:    utils.ToNullString(request.FirstName),
		LastName:     utils.ToNullString(request.LastName),
		DateOfBirth:  utils.ToNullTime(dob),
		Gender:       utils.ToNullString(request.Gender),
		HeightInches: utils.ToNullInt32(request.HeightInches),
	})

	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, profile, http.StatusCreated)
}
