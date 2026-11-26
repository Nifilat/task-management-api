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

const AdminRole = "admin"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

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
	claims, err := loadClaims(db, userID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateJWT(claims)
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

// loadClaims reads the user's current roles and permissions so the token always reflects the database.
func loadClaims(db *gorm.DB, userID uuid.UUID) (utils.Claims, error) {
	var user models.User
	if err := db.Preload("Roles.Permissions").First(&user, "id = ?", userID).Error; err != nil {
		return utils.Claims{}, err
	}

	claims := utils.Claims{UserID: userID.String(), Roles: []string{}, Permissions: []utils.PermissionClaim{}}
	actionsByResource := map[string]map[string]bool{}

	for _, role := range user.Roles {
		claims.Roles = append(claims.Roles, role.Name)
		if role.Name == AdminRole {
			claims.IsAdmin = true
		}
		for _, p := range role.Permissions {
			if actionsByResource[p.Resource] == nil {
				actionsByResource[p.Resource] = map[string]bool{}
			}
			actionsByResource[p.Resource][p.Action] = true
		}
	}

	for resource, actions := range actionsByResource {
		permission := utils.PermissionClaim{Resource: resource, Actions: []string{}}
		for action := range actions {
			permission.Actions = append(permission.Actions, action)
		}
		claims.Permissions = append(claims.Permissions, permission)
	}

	return claims, nil
}
