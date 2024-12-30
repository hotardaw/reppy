// TODO: add
// - rate limits for login
// - request body size limit to prevent mem exhaustion
// - request timeout contexts to prevent resource lockup/prevent DOS
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	requestTimeout = 10 * time.Second
	maxBodySize    = 1 * 1024 * 1024 // 1mb
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
	// Context times out after 10s
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel() // Always cancel to clean up resources

	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	if cleanPath != "/login" {
		response.SendError(w, "Invalid path", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" {
		response.SendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		response.SendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.auth.GenerateTokenPair(
		int64(user.UserID),
	)
	if err != nil {
		response.SendError(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// update last login
	if err := h.queries.UpdateLastLogin(ctx, user.UserID); err != nil {
		// log err but don't fail the request
		log.Printf("Failed to update last login: %v", err)
	}

	responseData := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SendSuccess(w, responseData)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	if cleanPath != "/refresh" {
		response.SendError(w, "Invalid path", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var request RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.RefreshToken) < 10 {
		response.SendError(w, "Invalid refresh token format", http.StatusBadRequest)
		return
	}

	claims, err := h.auth.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		response.SendError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.auth.GenerateTokenPair(
		claims.UserID,
	)
	if err != nil {
		response.SendError(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	responseData := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SendSuccess(w, responseData)
}
