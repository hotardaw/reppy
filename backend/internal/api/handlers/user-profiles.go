package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-reppy/backend/internal/api/response"
	"go-reppy/backend/internal/api/utils"
	"go-reppy/backend/internal/database/sqlc"
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
	switch r.Method {
	case http.MethodPost:
		h.CreateUserProfile(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *UserProfileHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("active") {
	case "true":
		h.GetAllActiveUserProfiles(w, r)
	case "false":
		h.GetAllInactiveUserProfiles(w, r)
	default:
		h.GetAllUserProfiles(w, r)
	}
}

// "/user-profiles"
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

	profile, err := h.queries.CreateUserProfile(r.Context(), sqlc.CreateUserProfileParams{
		UserID:       utils.ToNullInt32(request.UserID),
		FirstName:    utils.ToNullString(request.FirstName),
		LastName:     utils.ToNullString(request.LastName),
		DateOfBirth:  utils.ToNullTime(dob),
		Gender:       utils.ToNullString(request.Gender),
		HeightInches: utils.ToNullInt32(request.HeightInches),
	})
	if err != nil {
		response.SendError(w, "Failed to create user profile", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, profile, http.StatusCreated)
}

// "/user-profiles"
func (h *UserProfileHandler) GetAllUserProfiles(w http.ResponseWriter, r *http.Request) {
	userProfiles, err := h.queries.GetAllUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, "Failed to retrieve user profiles", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, userProfiles)
}

// "/user-profiles?active=true"
func (h *UserProfileHandler) GetAllActiveUserProfiles(w http.ResponseWriter, r *http.Request) {
	activeUserProfiles, err := h.queries.GetAllActiveUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, "Failed to retrieve active user profiles", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, activeUserProfiles)
}

// "/user-profiles?active=false"
func (h *UserProfileHandler) GetAllInactiveUserProfiles(w http.ResponseWriter, r *http.Request) {
	inactiveUserProfiles, err := h.queries.GetAllInactiveUserProfiles(r.Context())
	if err != nil {
		response.SendError(w, "Failed to retrieve inactive user profiles", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, inactiveUserProfiles)
}
