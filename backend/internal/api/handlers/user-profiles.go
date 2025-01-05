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

type CreateUserProfileRequest struct {
	UserID       int32  `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	HeightInches int32  `json:"height_inches"`
}

func (h *UserProfileHandler) HandleUserProfiles(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Handle both "/user-profiles" & "/user-profiles/active" endpoints
	if len(parts) != 2 || parts[1] != "user-profiles" {
		response.SendError(w, "Invalid URL - must be '/user-profiles'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("active") == "true" {
			h.GetAllActiveUserProfiles(w, r)
		} else if r.URL.Query().Get("active") == "false" {
			h.GetAllInactiveUserProfiles(w, r)
		} else {
			h.GetAllUserProfiles(w, r)
		}
	case http.MethodPost:
		h.CreateUserProfile(w, r)
	default:
		response.SendError(w, "Method not allowed - only GET and POST allowed at /user-profiles", http.StatusMethodNotAllowed)
		return
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

func (h *UserProfileHandler) GetAllActiveUserProfiles(w http.ResponseWriter, r *http.Request) {
	activeUserProfiles, err := h.queries.GetAllActiveUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, activeUserProfiles)
}

func (h *UserProfileHandler) GetAllInactiveUserProfiles(w http.ResponseWriter, r *http.Request) {
	inactiveUserProfiles, err := h.queries.GetAllInactiveUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, inactiveUserProfiles)
}

func (h *UserProfileHandler) CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	var request CreateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dob, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		response.SendError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	createUserProfileParams := sqlc.CreateUserProfileParams{
		UserID:       utils.ToNullInt32(request.UserID),
		FirstName:    utils.ToNullString(request.FirstName),
		LastName:     utils.ToNullString(request.LastName),
		DateOfBirth:  utils.ToNullTime(dob),
		Gender:       utils.ToNullString(request.Gender),
		HeightInches: utils.ToNullInt32(request.HeightInches),
	}
	profile, err := h.queries.CreateUserProfile(r.Context(), createUserProfileParams)
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, profile, http.StatusCreated)
}
