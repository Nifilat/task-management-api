package services

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// TODO(7): The starter signatures below only return an error. Change them to take the input and return the
// data, e.g. CreateTask(db, task) (models.Task, error) and GetTasks(db, ...) ([]models.Task, error).
// TODO(7): Implement them on TaskServiceImpl with GORM, then pass services.NewTaskService() in main.go.
type TaskService interface {
	CreateTask(db *gorm.DB) error
	UpdateTask(db *gorm.DB, taskID uuid.UUID) error
	DeleteTask(db *gorm.DB, taskID uuid.UUID) error
	GetTaskByID(db *gorm.DB, taskID uuid.UUID) error
	GetTasksByUser(db *gorm.DB, userID uuid.UUID) error
	GetTasks(db *gorm.DB) error
}
type TaskServiceImpl struct {
}

func NewTaskService() *TaskServiceImpl {
	return &TaskServiceImpl{}
}
