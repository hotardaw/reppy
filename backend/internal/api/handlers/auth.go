package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go-fitstat/backend/internal/api/middleware"
	"go-fitstat/backend/internal/api/response"
	"go-fitstat/backend/internal/database/sqlc"
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

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type SignupResponse struct {
	User         sqlc.User `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
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

// "/signup"
func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// tasks remaining here:
	// validate email format, password strength, username char requirements before processing
	// log failed signup attempts
	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Received request with method: %s", r.Method)
		return
	}

	// copying CreateUser API pretty closely here to avoid calling an API from within an API
	var request SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response.SendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	params := sqlc.CreateUserParams{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Username:     request.Username,
	}
	if params.Email == "" || params.PasswordHash == "" || params.Username == "" {
		response.SendError(w, "All fields must be filled", http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				response.SendError(w, "Email already in use", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				response.SendError(w, "Username already in use", http.StatusConflict)
			}
			return
		}
		response.SendError(w, "Duplicate value", http.StatusConflict)
		return
	}

	// generate token pair for user
	accessToken, refreshToken, err := h.auth.GenerateTokenPair(int64(user.UserID))
	if err != nil {
		response.SendError(w, "Created new user, but failed to generate tokens", http.StatusInternalServerError)
		return
	}

	responseData := SignupResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	// client extracts tokens from responseData to store in auth header
	// Authorization: Bearer <access_token>

	response.SendSuccess(w, responseData, http.StatusCreated)
}

// "/login"
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Received request with method: %s", r.Method)
		return
	}

	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" {
		response.SendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), request.Email)
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

	accessToken, refreshToken, err := h.auth.GenerateTokenPair(int64(user.UserID))
	if err != nil {
		response.SendError(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// update last login
	if err := h.queries.UpdateLastLogin(r.Context(), user.UserID); err != nil {
		// log err but don't fail the request
		log.Printf("Failed to update last login: %v", err)
	}

	responseData := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SendSuccess(w, responseData)
}

// "/refresh"
func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

// Since we're using JWTs, which can't be invalidated, we add the refresh token to the blacklist in the auth middleware
// in prod, we'd store the blacklist in redis or db instead of memory and implement blacklist cleanup
// in current impl, frontend needs to send a POST req to this logout endpoint with the refresh token stringify'ed in body, then localStorage.removeItem() on both tokens, clear React FE's auth state, and navigate to login or some unprotected page
// "/logout"
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// extract refresh token from client request
	var request RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.auth.InvalidateRefreshToken(request.RefreshToken); err != nil {
		response.SendError(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, nil)
}
