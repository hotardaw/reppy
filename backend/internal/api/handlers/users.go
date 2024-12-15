package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
	"strings"
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

	switch r.Method {
	case http.MethodGet:
		if len(parts) > 2 {
			h.GetUser(w, r, parts)
		} else {
			h.GetAllUsers(w, r, parts)
		}
	case http.MethodPost:
		if len(parts) > 2 {
			http.Error(w, "Invalid URL for POST request", http.StatusBadRequest)
			return
		}
		h.CreateUser(w, r)
	case http.MethodPut:
		if len(parts) > 2 {
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
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Username     string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
		Username:     request.Username,
	})

	if err != nil {
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
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Username     string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		UserID:       int32(id),
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
		Username:     request.Username,
	})

	if err != nil {
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
