package handlers

import (
	"errors"
	"net/http"
	"strings"

	"task-manager/backend/internal/models"
	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterHandler struct {
	db              *gorm.DB
	registerService services.RegisterService
}

func NewRegisterHandler(db *gorm.DB, registerService services.RegisterService) *RegisterHandler {
	return &RegisterHandler{db: db, registerService: registerService}
}

// RegisterInput is the request body for registration; it is separate from models.User so clients can't set ID or timestamps.
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Registration handles POST /api/v1/auth/register.
func (h *RegisterHandler) Registration(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, valid email and password are required"})
		return
	}

	username := strings.TrimSpace(input.Username)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be blank"})
		return
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: username, Email: input.Email, Password: input.Password}

	if err := h.registerService.RegisterUser(h.db, user); err != nil {
		switch {
		case errors.Is(err, services.ErrDuplicateUsername), errors.Is(err, services.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
