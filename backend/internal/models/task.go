package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Task is a unit of work owned by one user. It matches the tasks table (migration 000003).
type Task struct {
	gorm.Model
	ID          uuid.UUID  `json:"id" gorm:"primaryKey"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
}

// TaskInput is the request body for creating and updating a task, so clients cannot set ID or UserID.
type TaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
}
