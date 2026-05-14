package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"y/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLifetime = 7 * 24 * time.Hour
	bcryptCost    = 10
)

// ---------- secret / signing helpers --------------------------------------

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// Last-ditch default so the server still boots in dev. Logged loudly.
		log.Println("WARNING: JWT_SECRET not set, falling back to insecure default")
		return []byte("tally-dev-insecure-secret-change-me")
	}
	return []byte(s)
}

type AuthClaims struct {
	UserID   int    `json:"uid"`
	Username string `json:"u"`
	Role     string `json:"r"`
	jwt.RegisteredClaims
}

func signToken(u models.User) (string, error) {
	claims := AuthClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.Username,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret())
}

func parseToken(raw string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	t, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ---------- request-scoped user context -----------------------------------

type ctxKey string

const ctxKeyClaims ctxKey = "tally.claims"

// ClaimsFromContext returns the auth claims attached by RequireAuth, or nil
// if the request was anonymous.
func ClaimsFromContext(ctx context.Context) *AuthClaims {
	v := ctx.Value(ctxKeyClaims)
	if v == nil {
		return nil
	}
	c, _ := v.(*AuthClaims)
	return c
}

func withClaims(r *http.Request, c *AuthClaims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ctxKeyClaims, c))
}

// ---------- middlewares ---------------------------------------------------

// OptionalAuth populates the request context with claims when a valid token
// is present, but lets anonymous requests through.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := tryAuth(r); ok {
			r = withClaims(r, c)
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth rejects unauthenticated requests with 401.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := tryAuth(r)
		if !ok {
			httpError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(w, withClaims(r, c))
	})
}

// RequireAdmin requires a valid token AND the admin role.
func RequireAdmin(next http.Handler) http.Handler {
	return RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := ClaimsFromContext(r.Context())
		if c == nil || c.Role != string(models.RoleAdmin) {
			httpError(w, http.StatusForbidden, "admin role required")
			return
		}
		next.ServeHTTP(w, r)
	}))
}

func tryAuth(r *http.Request) (*AuthClaims, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return nil, false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, false
	}
	c, err := parseToken(parts[1])
	if err != nil {
		return nil, false
	}
	return c, true
}

// ---------- handlers ------------------------------------------------------

var (
	reUsername = regexp.MustCompile(`^[A-Za-z0-9_]{3,30}$`)
	reEmail    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// Signup creates a new user with role "user" and returns a fresh JWT.
func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	switch {
	case !reUsername.MatchString(req.Username):
		httpError(w, http.StatusBadRequest, "username must be 3-30 chars, letters/digits/underscore")
		return
	case !reEmail.MatchString(req.Email):
		httpError(w, http.StatusBadRequest, "invalid email")
		return
	case len(req.Password) < 6:
		httpError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		log.Printf("Signup hash: %v", err)
		httpError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	var u models.User
	err = DB().QueryRow(
		`INSERT INTO users (username, email, password_hash, role)
		 VALUES ($1, $2, $3, 'user')
		 RETURNING id, username, email, role, created_at`,
		req.Username, req.Email, string(hash),
	).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			httpError(w, http.StatusConflict, "username or email already exists")
			return
		}
		log.Printf("Signup insert: %v", err)
		httpError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	token, err := signToken(u)
	if err != nil {
		log.Printf("Signup sign: %v", err)
		httpError(w, http.StatusInternalServerError, "could not issue token")
		return
	}
	writeJSON(w, http.StatusCreated, models.AuthResponse{Token: token, User: u})
}

// Login looks up a user by username OR email and verifies their password.
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.UsernameOrEmail = strings.TrimSpace(req.UsernameOrEmail)
	if req.UsernameOrEmail == "" || req.Password == "" {
		httpError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	var (
		u    models.User
		hash string
	)
	err := DB().QueryRow(
		`SELECT id, username, email, password_hash, role, created_at
		 FROM users
		 WHERE username = $1 OR email = LOWER($1)`,
		req.UsernameOrEmail,
	).Scan(&u.ID, &u.Username, &u.Email, &hash, &u.Role, &u.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		httpError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		log.Printf("Login query: %v", err)
		httpError(w, http.StatusInternalServerError, "login failed")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		httpError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := signToken(u)
	if err != nil {
		log.Printf("Login sign: %v", err)
		httpError(w, http.StatusInternalServerError, "could not issue token")
		return
	}
	writeJSON(w, http.StatusOK, models.AuthResponse{Token: token, User: u})
}

// Me returns the authenticated user's profile.
func Me(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var u models.User
	err := DB().QueryRow(
		`SELECT id, username, email, role, created_at FROM users WHERE id=$1`,
		c.UserID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		httpError(w, http.StatusUnauthorized, "user no longer exists")
		return
	}
	if err != nil {
		log.Printf("Me query: %v", err)
		httpError(w, http.StatusInternalServerError, "could not load user")
		return
	}
	writeJSON(w, http.StatusOK, u)
}
