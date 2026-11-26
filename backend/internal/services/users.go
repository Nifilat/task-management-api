package services

import (
	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserProfile(db *gorm.DB, userID uuid.UUID) (models.User, error)
	GetUsers(db *gorm.DB) ([]models.User, error)
	DeleteUser(db *gorm.DB, userId uuid.UUID) error
}

type UserServiceImpl struct {
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (s *UserServiceImpl) GetUserProfile(db *gorm.DB, userID uuid.UUID) (models.User, error) {
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetUsers(db *gorm.DB) ([]models.User, error) {
	users := []models.User{}

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserServiceImpl) DeleteUser(db *gorm.DB, userId uuid.UUID) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.User{}, "id = ?", userId)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		// Remove their refresh tokens so a deleted user cannot get new access tokens.
		return tx.Unscoped().Delete(&models.Token{}, "user_id = ?", userId).Error
	})
}
