package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// User maps to the users table.
// TODO(5): Add the user's roles (many2many user_roles) once the RBAC tables exist.
type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}
