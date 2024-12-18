// For /users endpoint - no specific users
package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Ensure only /users endpoint is handled
	if len(parts) != 2 || parts[1] != "users" {
		http.Error(w, "Invalid URL - must be '/users'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r, parts)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Username:     request.Username,
	})

	if err != nil {
		// handle dupe emails/usernames
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				http.Error(w, "Email already in use", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				http.Error(w, "Username already taken", http.StatusConflict)
			} else {
				http.Error(w, "Duplicate value", http.StatusConflict)
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/*
1. Define handler funcs that process incoming HTTP reqs
2. Map URLs/routes to their corresponding handlers
3. Handle req validation, auth, and authorization
4. Format & send HTTP responses

Import the SQLc queries
Define the HTTP hanndler functions that use these queries
Handle the HTTP request/response lifecycle

Then register these handlers to make them accessible
*/
