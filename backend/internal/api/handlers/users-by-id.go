package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
)

type UserByIDHandler struct {
	queries *sqlc.Queries
}

func NewUserByIDHandler(q *sqlc.Queries) *UserByIDHandler {
	return &UserByIDHandler{
		queries: q,
	}
}

type UpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Username string `json:"username,omitempty"`
}

func (h *UserByIDHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid user ID in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, int32(userID))
	case http.MethodPatch:
		h.UpdateUser(w, r, int32(userID))
	case http.MethodDelete:
		h.DeleteUser(w, r, int32(userID))
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/users/{id}"
func (h *UserByIDHandler) GetUser(w http.ResponseWriter, r *http.Request, userID int32) {
	user, err := h.queries.GetUser(r.Context(), userID)
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user)
}

// use utils.isValidXYZ to validate email, un, pw
// "/users/{id}"
func (h *UserByIDHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userID int32) {
	var request UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "" && request.Password == "" && request.Username == "" {
		response.SendError(w, "No fields to update", http.StatusBadRequest)
		return
	}

	updateUserParams := sqlc.UpdateUserParams{
		UserID:   userID,
		Email:    request.Email,
		Username: request.Username,
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			response.SendError(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		updateUserParams.PasswordHash = string(hashedPassword)
	}

	user, err := h.queries.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				response.SendError(w, "Email already in use", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				response.SendError(w, "Username already taken", http.StatusConflict)
			} else {
				response.SendError(w, "Duplicate value", http.StatusConflict)
			}
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user)
}

// "/users/{id}"
func (h *UserByIDHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userID int32) {
	deletedUser, err := h.queries.DeleteUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message": "Muscle deleted successfully",
		"id":      deletedUser,
	}, http.StatusOK) // Not StatusNoContent bc this is a soft delete
}

/*
map[string]interface{}{
		"message": "Muscle deleted successfully",
		"id":      deletedMuscle,
	}
*/

func parseUserID(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, errors.New("invalid path format")
	}
	return strconv.Atoi(parts[2])
}
