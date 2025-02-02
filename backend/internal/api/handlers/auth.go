package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"go-reppy/backend/internal/api/middleware"
	"go-reppy/backend/internal/api/response"
	"go-reppy/backend/internal/api/utils"
	"go-reppy/backend/internal/database/sqlc"
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

type GoogleAuthHandler struct {
	oauth2Config *oauth2.Config
	queries      *sqlc.Queries
	auth         *middleware.AuthMiddleware
}

func NewGoogleAuthHandler(queries *sqlc.Queries, auth *middleware.AuthMiddleware, config middleware.JWTConfig) *GoogleAuthHandler {
	oauth2Config := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleAuthHandler{
		oauth2Config: oauth2Config,
		queries:      queries,
		auth:         auth,
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

/* Google login flow:
Frontend (localhost:8080)
  → Google Sign-in
  → Backend Callback (localhost:8081/auth/google/callback)
  → Frontend (localhost:8080) with tokens
*/

func (h *GoogleAuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// to prevent CSRF - gen random string that gets validated when Google redirects user back here
	state := generateRandomState()

	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	url := h.oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// endpoint user is redirected to by Google after sign-in success
func (h *GoogleAuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// get state from cookie
	cookie, err := r.Cookie("oauthstate")
	if err != nil {
		response.SendError(w, "State cookie not found", http.StatusBadRequest)
		return
	}
	// verify state match
	if r.FormValue("state") != cookie.Value {
		response.SendError(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// exchange auth code for a token
	code := r.URL.Query().Get("code")
	token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Google auth failed: %v", err)
		response.SendError(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// get user info from google
	client := h.oauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get Google user info: %v", err)
		response.SendError(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Printf("Failed to decode Google user info: %v", err)
		response.SendError(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}
	if !googleUser.VerifiedEmail {
		response.SendError(w, "Google email must be verified", http.StatusBadRequest)
		return
	}

	// check if user exists
	user, err := h.queries.GetUserByGoogleID(r.Context(), utils.ToNullString(googleUser.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// check if email exists in local auth
			existingUser, err := h.queries.GetUserByEmail(r.Context(), googleUser.Email)
			if err == nil && existingUser.AuthProvider == utils.ToNullString("local") {
				response.SendError(w, "Email already registered with password", http.StatusConflict)
				return
			}

			// create new reppy user
			user, err = h.queries.CreateGoogleUser(r.Context(), sqlc.CreateGoogleUserParams{
				Email:    googleUser.Email,
				Username: googleUser.Name,
				GoogleID: utils.ToNullString(googleUser.ID),
			})
			if err != nil {
				log.Printf("Failed to create Google user: %v", err)
				if strings.Contains(err.Error(), "unique constraint") {
					if strings.Contains(err.Error(), "email") {
						response.SendError(w, "Email already in use", http.StatusConflict)
						return
					} else if strings.Contains(err.Error(), "username") {
						response.SendError(w, "Username already in use", http.StatusConflict)
						return
					}
				}

				response.SendError(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Database error: %v", err)
			response.SendError(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// grant JWTs
	accessToken, refreshToken, err := h.auth.GenerateTokenPair(int64(user.UserID))
	if err != nil {
		log.Printf("Failed to generate tokens: %v", err)
		response.SendError(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// success!
	// frontendURL := "http://localhost:8080"
	frontendURL := "http://localhost:8081"
	http.Redirect(w, r, fmt.Sprintf("%s?access_token=%s&refresh_token=%s",
		frontendURL, accessToken, refreshToken), http.StatusTemporaryRedirect)

}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
