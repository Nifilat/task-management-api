package services

import (
	"errors"
	"time"

	"task-manager/backend/internal/models"
	"task-manager/backend/internal/utils"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// TODO(5): Add a method that loads a user's roles and permissions for the token claims.
type AuthService interface {
	LoginUser(db *gorm.DB, username, password string) (*models.User, error)
	GenerateToken(db *gorm.DB, userID uuid.UUID) (string, string, error)
	RotateRefreshToken(db *gorm.DB, refreshToken string) (string, string, error)
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

func (s *AuthServiceImpl) LoginUser(db *gorm.DB, username, password string) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !VerifyPassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}
	return &user, nil
}

// GenerateToken returns a new access_token (JWT) and refresh_token, and stores the refresh_token in the tokens table.
func (s *AuthServiceImpl) GenerateToken(db *gorm.DB, userID uuid.UUID) (string, string, error) {
	accessToken, err := utils.GenerateJWT(utils.Claims{UserID: userID.String()})
	if err != nil {
		return "", "", err
	}

	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}
	refreshToken, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}

	token := models.Token{
		ID:           tokenID,
		UserId:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(utils.AccessTokenTTL),
	}
	if err := db.Create(&token).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.String(), nil
}

// RotateRefreshToken exchanges a valid refresh_token for a new access/refresh pair; the old one is deleted so it can't be reused.
func (s *AuthServiceImpl) RotateRefreshToken(db *gorm.DB, refreshToken string) (string, string, error) {
	refreshID, err := uuid.FromString(refreshToken)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	var access, refresh string
	err = db.Transaction(func(tx *gorm.DB) error {
		var stored models.Token
		if err := tx.Where("refresh_token = ?", refreshID).First(&stored).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		if err := tx.Unscoped().Delete(&models.Token{}, "id = ?", stored.ID).Error; err != nil {
			return err
		}
		if time.Now().After(stored.ExpiresAt) {
			return ErrInvalidRefreshToken
		}

		access, refresh, err = s.GenerateToken(tx, stored.UserId)
		return err
	})
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
