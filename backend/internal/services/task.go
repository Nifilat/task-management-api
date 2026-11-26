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
	GetTasksByUser(db *gorm.DB, userID uuid.UUID, opts TaskListOptions) ([]models.Task, int64, error)
	GetTasks(db *gorm.DB, opts TaskListOptions) ([]models.Task, int64, error)
}

// TaskListOptions holds validated paging and sorting values (SortBy is a key of TaskSortFields).
type TaskListOptions struct {
	Limit  int
	Offset int
	SortBy string
	Order  string
}

// TaskSortFields whitelists the sortable fields; the value is the column that is safe to put in ORDER BY.
var TaskSortFields = map[string]string{
	"created_at": "created_at",
	"due_date":   "due_date",
	"title":      "title",
	"status":     "status",
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

func (s *TaskServiceImpl) GetTasksByUser(db *gorm.DB, userID uuid.UUID, opts TaskListOptions) ([]models.Task, int64, error) {
	return listTasks(db.Where("user_id = ?", userID), opts)
}

func (s *TaskServiceImpl) GetTasks(db *gorm.DB, opts TaskListOptions) ([]models.Task, int64, error) {
	return listTasks(db, opts)
}

// listTasks counts the matching tasks and returns one sorted page of them.
func listTasks(query *gorm.DB, opts TaskListOptions) ([]models.Task, int64, error) {
	var total int64
	if err := query.Model(&models.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	tasks := []models.Task{}
	err := query.
		Order(TaskSortFields[opts.SortBy] + " " + opts.Order).
		Order("id").
		Limit(opts.Limit).
		Offset(opts.Offset).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}
