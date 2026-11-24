package services

import (
	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TODO(2): Implement LoginUser and GenerateToken on AuthServiceImpl (then the interface is satisfied and
// main.go can pass services.NewAuthService() instead of nil).
// TODO(2): Add a method to find the stored refresh_token, and one to save a new one (tokens table).
// TODO(5): Add a method that loads a user's roles and permissions for the token claims.
type AuthService interface {
	LoginUser(db *gorm.DB, username, password string) (*models.User, error)
	GenerateToken(db *gorm.DB, userID uuid.UUID) (string, string, error)
}

type AuthServiceImpl struct {
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
