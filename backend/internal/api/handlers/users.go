// For /users endpoint - no specific users
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
)

// Holds dependencies for user handlers
type UserHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(q *sqlc.Queries) *UserHandler {
	return &UserHandler{
		queries: q,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response.SendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	createUserParams := sqlc.CreateUserParams{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Username:     request.Username,
	}
	if createUserParams.Email == "" || createUserParams.PasswordHash == "" || createUserParams.Username == "" {
		response.SendError(w, "All fields must be filled", http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				response.SendError(w, "Email already in use", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				response.SendError(w, "Username already in use", http.StatusConflict)
			}
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user, http.StatusCreated)
}
