package services

import (
	"task-manager/backend/internal/models"

	"gorm.io/gorm"
)

// TODO(1): Implement RegisterUser on RegisterServiceImpl: hash with bcrypt, set the UUID, create the row.
// TODO(5): In the same transaction, assign the default "user" role.
type RegisterService interface {
	RegisterUser(db *gorm.DB, user models.User) error
}

type RegisterServiceImpl struct{}

func NewRegisterService() *RegisterServiceImpl {
	return &RegisterServiceImpl{}
}
