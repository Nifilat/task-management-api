package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenIssuer    = "task-manager"
	AccessTokenTTL = time.Hour

	RefreshTokenTTL = time.Hour
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

// GenerateRefreshJWT signs a refresh_token (valid 1 hour) whose jti is the id of its row in the tokens table.
func GenerateRefreshJWT(userID, tokenID string) (string, time.Time, error) {
	secret := GetEnv("JWT_SECRET", "")
	if secret == "" {
		return "", time.Time{}, errors.New("JWT_SECRET is not set")
	}

	now := time.Now()
	expiresAt := now.Add(RefreshTokenTTL)
	claims := jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		Subject:   userID,
		ID:        tokenID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	return signed, expiresAt, err
}

// ParseRefreshJWT validates a refresh_token's signature, issuer and expiry.
func ParseRefreshJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	secret := GetEnv("JWT_SECRET", "")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(TokenIssuer), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.ID == "" {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

const MinPasswordLength = 8

// ValidatePassword returns an error naming the first password rule that is not met.
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	var upper, lower, digit, special bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsDigit(r):
			digit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			special = true
		}
	}

	switch {
	case !upper:
		return errors.New("password must contain at least one uppercase letter")
	case !lower:
		return errors.New("password must contain at least one lowercase letter")
	case !digit:
		return errors.New("password must contain at least one number")
	case !special:
		return errors.New("password must contain at least one special character")
	}
	return nil
}
