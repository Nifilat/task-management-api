package services

import (
	"errors"

	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)

const DefaultRole = "user"

type RegisterService interface {
	RegisterUser(db *gorm.DB, user models.User) error
}

type RegisterServiceImpl struct{}

func NewRegisterService() *RegisterServiceImpl {
	return &RegisterServiceImpl{}
}

func (s *RegisterServiceImpl) RegisterUser(db *gorm.DB, user models.User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrDuplicateUsername
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)

		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		user.ID = id

		var role models.Role
		if err := tx.Where("name = ?", DefaultRole).First(&role).Error; err != nil {
			return err
		}
		user.Roles = []models.Role{role}

		if err := tx.Omit("Roles.*").Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicateEmail
			}
			return err
		}
		return nil
	})
}
