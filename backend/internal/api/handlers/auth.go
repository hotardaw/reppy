package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/database/sqlc"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *sqlc.Queries
	auth    *middleware.AuthMiddleware
}

func NewAuthHandler(queries *sqlc.Queries, auth *middleware.AuthMiddleware) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		auth:    auth,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get user from db
	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// generate tokens
	accessToken, refreshToken, err := h.auth.GenerateTokenPair(
		int64(user.UserID),
		user.Email,
		user.Username,
	)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// update last login
	if err := h.queries.UpdateLastLogin(r.Context(), sqlc.UpdateLastLoginParams{
		UserID: user.UserID,
		LastLogin: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}); err != nil {
		// log err but don't fail the request
		log.Printf("Failed to update last login: %v", err)
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	/*
		check if POST method
		parse JSON req body
		get user from db via email
		compare password hash
		gen token pair
		update last login
		return tokens as JSON resp
	*/
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate refresh token
	claims, err := h.auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// gen new token pair
	accessToken, refreshToken, err := h.auth.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Username,
	)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*
In main.go:
- Create auth middleware instance with config
- Create auth handler instance
- Add login and refresh routes to mux

SQL queries needed:
- Get user by email
- Update last login timestamp
*/
