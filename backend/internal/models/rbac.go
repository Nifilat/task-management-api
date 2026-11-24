package models

import "github.com/gofrs/uuid"

// Role is a named group of permissions ("user", "admin"). Maps to the roles table.
// TODO(4): Add GORM tags/relations (many2many user_roles) once the migrations are written.
type Role struct {
	ID   uuid.UUID `json:"id" gorm:"primaryKey"`
	Name string    `json:"name" gorm:"unique"`
}

// Permission is an action on a resource, e.g. resource "task", action "delete". Maps to the permissions table.
// TODO(4): Add a many2many role_permissions relation on Role.
type Permission struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Resource string    `json:"resource"`
	Action   string    `json:"action"`
}
