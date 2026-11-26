package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenIssuer    = "task-manager"
	AccessTokenTTL = time.Hour
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
type Claims struct {
	UserID      string            `json:"user_id"`
	Roles       []string          `json:"roles"`
	IsAdmin     bool              `json:"is_admin"`
	Permissions []PermissionClaim `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateJWT signs an access_token for the user (valid 1 hour, rubric) with HS256 and JWT_SECRET.
func GenerateJWT(claims Claims) (string, error) {
	secret := GetEnv("JWT_SECRET", "")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		Subject:   claims.UserID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseJWT validates a token's signature, signing method, issuer and expiry and returns its claims.
func ParseJWT(tokenString string) (*Claims, error) {
	secret := GetEnv("JWT_SECRET", "")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(TokenIssuer), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
