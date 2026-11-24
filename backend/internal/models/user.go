package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// User maps to the users table.
// TODO(1): SECURITY: Password is serialized into JSON responses as-is. Tag it `json:"-"` so the hash is never
// returned, and use a separate request struct for registration (a `json:"-"` field can't be bound from input).
// TODO(5): Add the user's roles (many2many user_roles) once the RBAC tables exist.
type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Username string    `json:"username" gorm:"unique"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
}
