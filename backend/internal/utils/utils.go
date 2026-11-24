package utils

import (
	"errors"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// ---------------------------------------------------------------------------
// Shared helpers (rubric stand-out: keep JWT parsing and UUID validation here, DRY).
// ---------------------------------------------------------------------------

// PermissionClaim lists the actions a token holder may perform on one resource.
type PermissionClaim struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

// Claims is the JWT payload. The frontend decodes user_id, roles, is_admin and permissions from it.
// TODO(5): Embed jwt.RegisteredClaims (from github.com/golang-jwt/jwt/v5) for exp, iss, sub.
type Claims struct {
	UserID      string            `json:"user_id"`
	Roles       []string          `json:"roles"`
	IsAdmin     bool              `json:"is_admin"`
	Permissions []PermissionClaim `json:"permissions"`
}

// GenerateJWT signs an access_token for the user (valid 1 hour, rubric).
// TODO(2): Build the claims, sign with HS256 using the JWT_SECRET env var, return the token string.
// TODO(5): Include the user's roles and permissions in the claims.
func GenerateJWT(claims Claims) (string, error) {
	return "", errors.New("not implemented")
}

// ParseJWT validates a token's signature and expiry and returns its claims.
// TODO(6): Verify the signing method is HS256 (reject "none"), the issuer, and expiry; return the claims.
func ParseJWT(tokenString string) (*Claims, error) {
	return nil, errors.New("not implemented")
}

// IsValidUUID reports whether s is a well-formed UUID (use before touching the database).
// TODO(8): Implement with uuid.FromString from github.com/gofrs/uuid.
func IsValidUUID(s string) bool {
	return false
}
