package handlers

import (
	"errors"
	"net/http"

	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db          *gorm.DB
	taskService services.TaskService
}

func NewTaskHandler(db *gorm.DB, taskService services.TaskService) *TaskHandler {
	return &TaskHandler{db: db, taskService: taskService}
}

// CreateTask handles POST /api/v1/tasks.
// TODO(7): Bind title, description, due_date; 400 on invalid input. Owner = the user_id from the token claims.
// TODO(7): Generate the UUID, save with GORM, return 201 + the created task.
func (h *TaskHandler) CreateTask(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{"message": "task created successfully"})
}

// UpdateTask handles PUT /api/v1/tasks/:id.
// TODO(7): Validate the :id UUID (utils.IsValidUUID); 404 if the task doesn't exist.
// TODO(7): ABAC: a user may update only their own task; an admin may update any (403 otherwise).
func (h *TaskHandler) UpdateTask(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask handles DELETE /api/v1/tasks/:id.
// TODO(7): 404 if missing. Permissions table: only admins can delete tasks (task:delete).
func (h *TaskHandler) DeleteTask(c *gin.Context) {

	c.JSON(http.StatusNoContent, nil)
}

// GetTaskByID handles GET /api/v1/tasks/:id.
// TODO(7): 404 if missing. ABAC: users see only their own task, admins any.
func (h *TaskHandler) GetTaskByID(c *gin.Context) {

	c.JSON(http.StatusOK, "task by id")
}

// GetTasksByUser handles GET /api/v1/users/:user_id/tasks (the frontend's task list calls this).
// TODO(7): A user may list only their own tasks (:user_id must equal the token's user_id); admins any user's.
func (h *TaskHandler) GetTasksByUser(c *gin.Context) {

	c.JSON(http.StatusOK, "tasks list by user")
}

// GetTasks handles GET /api/v1/tasks.
// TODO(7): Admins get every task; a regular user gets only their own.
// TODO(9, optional): pagination (limit/offset) and sorting, with a whitelisted sort field.
func (h *TaskHandler) GetTasks(c *gin.Context) {

	c.JSON(http.StatusOK, "tasks list")
}

func handleTaskError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process task request",
		})
	}
}
