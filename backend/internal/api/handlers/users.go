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
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/")) // "/users"
	parts := strings.Split(cleanPath, "/")                       // "[ users ]"

	switch r.Method {
	case http.MethodGet:
		if len(parts) > 2 {
			h.GetUser(w, r, parts)
		} else {
			h.GetAllUsers(w, r, parts)
		}
	case http.MethodPost:
		if len(parts) > 3 {
			http.Error(w, "Invalid URL for POST request", http.StatusBadRequest)
			return
		}
		h.CreateUser(w, r)
	case http.MethodPut:
		if len(parts) > 3 {
			http.Error(w, "Invalid URL for PUT request", http.StatusBadRequest)
			return
		}
		h.UpdateUser(w, r, parts)
	case http.MethodDelete:
		if len(parts) > 2 {
			http.Error(w, "Invalid URL for DELETE request", http.StatusBadRequest)
			return
		}
		h.DeleteUser(w, r, parts)
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUser(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
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

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateUserParams{
		UserID:   int32(id),
		Email:    request.Email,
		Username: request.Username,
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		params.PasswordHash = string(hashedPassword)
	}

	user, err := h.queries.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		// handle dupes
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteUser(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
