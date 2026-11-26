package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// User maps to the users table.
type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Roles    []Role    `json:"roles,omitempty" gorm:"many2many:user_roles"`
}
