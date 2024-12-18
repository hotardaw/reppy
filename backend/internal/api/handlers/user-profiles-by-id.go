// For /user-profiles/{user_id} endpoint
package handlers

import (
	"database/sql"
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
	"strings"
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
	case http.MethodPut:
		h.UpdateUserProfile(w, r, parts)
	case http.MethodDelete:
		h.DeleteUserProfile(w, r, parts)
	default:
		http.Error(w, "Method not allowed - only GET, PUT, and DELETE", http.StatusMethodNotAllowed)
	}
}

func (h *UserProfileByIDHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
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

func (h *UserProfileByIDHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
}
func (h *UserProfileByIDHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request, parts []string) {
}
