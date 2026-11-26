package models

import "github.com/gofrs/uuid"

// Role is a named group of permissions ("user", "admin"). Maps to the roles table.
type Role struct {
	ID          uuid.UUID    `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions"`
}

// Permission is an action on a resource, e.g. resource "task", action "delete". Maps to the permissions table.
type Permission struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Resource string    `json:"resource"`
	Action   string    `json:"action"`
}
