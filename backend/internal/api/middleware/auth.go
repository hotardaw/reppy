// JWT verification, session checking
package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"go-fitsync/backend/internal/api/response"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrMissingAuth  = errors.New("missing authorization header")
	ErrInvalidAuth  = errors.New("invalid authorization format")
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

type JWTConfig struct {
	AccessSecret    []byte
	RefreshSecret   []byte
	AccessDuration  time.Duration // 15m
	RefreshDuration time.Duration // 7d
	Issuer          string
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"` // either "access" or "refresh"
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	config            JWTConfig
	blacklistedTokens map[string]time.Time
}

func NewAuthMiddleware(config JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{
		config:            config,
		blacklistedTokens: make(map[string]time.Time),
	}
}

func (am *AuthMiddleware) InvalidateRefreshToken(token string) error {
	am.blacklistedTokens[token] = time.Now()
	return nil
}

// generate access & refresh token
func (am *AuthMiddleware) GenerateTokenPair(userID int64) (accessToken, refreshToken string, err error) {
	accessToken, err = am.generateToken(userID, "access", am.config.AccessSecret, am.config.AccessDuration)
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err = am.generateToken(userID, "refresh", am.config.RefreshSecret, am.config.RefreshDuration)
	if err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (am *AuthMiddleware) generateToken(userID int64, tokenType string, secret []byte, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType, // access or refresh
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // "exp"
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // "iat"
			NotBefore: jwt.NewNumericDate(time.Now()),               // "nbf"
			Issuer:    am.config.Issuer,                             // "iss"
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (am *AuthMiddleware) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return am.validateToken(tokenString, am.config.RefreshSecret)
}

func (am *AuthMiddleware) validateToken(tokenString string, secret []byte) (*Claims, error) {
	// check if token is blacklisted first
	if _, blacklisted := am.blacklistedTokens[tokenString]; blacklisted {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (am *AuthMiddleware) extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuth
	}

	parts := strings.Split(authHeader, " ")
	const AuthSchemeBearer = "bearer"

	if len(parts) != 2 || strings.ToLower(parts[0]) != AuthSchemeBearer {
		return "", ErrInvalidAuth
	}

	return parts[1], nil
}

// for protected routes
func (am *AuthMiddleware) AuthenticateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := am.extractTokenFromHeader(r)
		if err != nil {
			response.SendError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := am.validateToken(tokenString, am.config.AccessSecret)
		if err != nil {
			if err == ErrExpiredToken {
				response.SendError(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			response.SendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.TokenType != "access" {
			response.SendError(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	claims, ok := ctx.Value(UserClaimsKey).(*Claims)
	if !ok {
		return 0, errors.New("no claims in context")
	}
	return claims.UserID, nil
}
