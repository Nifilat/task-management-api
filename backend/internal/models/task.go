package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Task is a unit of work owned by one user. It must match the tasks table (migration 000003).
// TODO(7): Confirm the fields against your migration; add Status validation (pending|in_progress|completed).
// TODO(7): Use a request struct (title, description, due_date) for create/update so clients cannot set ID or UserID.
type Task struct {
	gorm.Model
	ID          uuid.UUID  `json:"id" gorm:"primaryKey"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
}
