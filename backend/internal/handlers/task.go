package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-manager/backend/internal/models"
	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db          *gorm.DB
	taskService services.TaskService
}

func NewTaskHandler(db *gorm.DB, taskService services.TaskService) *TaskHandler {
	return &TaskHandler{db: db, taskService: taskService}
}

// callerFromContext returns the token claims and the caller's user id set by AuthMiddleware.
func callerFromContext(c *gin.Context) (*utils.Claims, uuid.UUID, bool) {
	value, exists := c.Get("claims")
	claims, ok := value.(*utils.Claims)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return nil, uuid.Nil, false
	}

	userID, err := uuid.FromString(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return nil, uuid.Nil, false
	}
	return claims, userID, true
}

// loadTask finds the task in :id and checks the caller may access it: the owner, or an admin (ABAC).
func (h *TaskHandler) loadTask(c *gin.Context) (models.Task, bool) {
	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return models.Task{}, false
	}

	claims, userID, ok := callerFromContext(c)
	if !ok {
		return models.Task{}, false
	}

	task, err := h.taskService.GetTaskByID(h.db, taskID)
	if err != nil {
		handleTaskError(c, err)
		return models.Task{}, false
	}

	if !claims.IsAdmin && task.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only access your own tasks"})
		return models.Task{}, false
	}
	return task, true
}

// CreateTask handles POST /api/v1/tasks. The owner is the user in the token, never the request body.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required; status must be pending, in_progress or completed"})
		return
	}

	_, userID, ok := callerFromContext(c)
	if !ok {
		return
	}

	task := models.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Status:      input.Status,
	}
	if task.Status == "" {
		task.Status = "pending"
	}

	created, err := h.taskService.CreateTask(h.db, task)
	if err != nil {
		handleTaskError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateTask handles PUT /api/v1/tasks/:id (owner or admin).
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	task, ok := h.loadTask(c)
	if !ok {
		return
	}

	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required; status must be pending, in_progress or completed"})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.DueDate = input.DueDate
	if input.Status != "" {
		task.Status = input.Status
	}

	updated, err := h.taskService.UpdateTask(h.db, task)
	if err != nil {
		handleTaskError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteTask handles DELETE /api/v1/tasks/:id. The route requires task:delete, which only admins have.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.taskService.DeleteTask(h.db, taskID); err != nil {
		handleTaskError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetTaskByID handles GET /api/v1/tasks/:id (owner or admin).
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	task, ok := h.loadTask(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, task)
}

// GetTasksByUser handles GET /api/v1/users/:user_id/tasks (the frontend's task list calls this).
// The route's RequireOwnerOrAdmin middleware already ensures the caller is that user or an admin.
func (h *TaskHandler) GetTasksByUser(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	opts, ok := parseTaskListOptions(c)
	if !ok {
		return
	}

	tasks, total, err := h.taskService.GetTasksByUser(h.db, userID, opts)
	if err != nil {
		handleTaskError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, tasks)
}

// GetTasks handles GET /api/v1/tasks: admins get every task, a regular user gets only their own.
// Query params (all optional): limit (1-100, default 100), offset, sort_by, order.
func (h *TaskHandler) GetTasks(c *gin.Context) {
	claims, userID, ok := callerFromContext(c)
	if !ok {
		return
	}

	opts, ok := parseTaskListOptions(c)
	if !ok {
		return
	}

	var tasks []models.Task
	var total int64
	var err error
	if claims.IsAdmin {
		tasks, total, err = h.taskService.GetTasks(h.db, opts)
	} else {
		tasks, total, err = h.taskService.GetTasksByUser(h.db, userID, opts)
	}
	if err != nil {
		handleTaskError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, tasks)
}

// parseTaskListOptions reads and validates limit, offset, sort_by and order; it replies 400 and returns false if any is invalid.
func parseTaskListOptions(c *gin.Context) (services.TaskListOptions, bool) {
	const maxLimit = 100
	opts := services.TaskListOptions{Limit: maxLimit, SortBy: "created_at", Order: "asc"}

	if v := c.Query("limit"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil || limit < 1 || limit > maxLimit {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number between 1 and 100"})
			return opts, false
		}
		opts.Limit = limit
	}

	if v := c.Query("offset"); v != "" {
		offset, err := strconv.Atoi(v)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be a number that is 0 or more"})
			return opts, false
		}
		opts.Offset = offset
	}

	if v := c.Query("sort_by"); v != "" {
		if _, ok := services.TaskSortFields[v]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sort_by must be created_at, due_date, title or status"})
			return opts, false
		}
		opts.SortBy = v
	}

	if v := strings.ToLower(c.Query("order")); v != "" {
		if v != "asc" && v != "desc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order must be asc or desc"})
			return opts, false
		}
		opts.Order = v
	}

	return opts, true
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
