package services

import (
	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(db *gorm.DB, task models.Task) (models.Task, error)
	UpdateTask(db *gorm.DB, task models.Task) (models.Task, error)
	DeleteTask(db *gorm.DB, taskID uuid.UUID) error
	GetTaskByID(db *gorm.DB, taskID uuid.UUID) (models.Task, error)
	GetTasksByUser(db *gorm.DB, userID uuid.UUID) ([]models.Task, error)
	GetTasks(db *gorm.DB) ([]models.Task, error)
}
type TaskServiceImpl struct {
}

func NewTaskService() *TaskServiceImpl {
	return &TaskServiceImpl{}
}

func (s *TaskServiceImpl) CreateTask(db *gorm.DB, task models.Task) (models.Task, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return models.Task{}, err
	}
	task.ID = id

	if err := db.Create(&task).Error; err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *TaskServiceImpl) UpdateTask(db *gorm.DB, task models.Task) (models.Task, error) {
	if err := db.Save(&task).Error; err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *TaskServiceImpl) DeleteTask(db *gorm.DB, taskID uuid.UUID) error {
	result := db.Delete(&models.Task{}, "id = ?", taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *TaskServiceImpl) GetTaskByID(db *gorm.DB, taskID uuid.UUID) (models.Task, error) {
	var task models.Task
	if err := db.First(&task, "id = ?", taskID).Error; err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetTasksByUser(db *gorm.DB, userID uuid.UUID) ([]models.Task, error) {
	tasks := []models.Task{}
	if err := db.Where("user_id = ?", userID).Order("created_at").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetTasks(db *gorm.DB) ([]models.Task, error) {
	tasks := []models.Task{}
	if err := db.Order("created_at").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
