package handlers

import (
	"errors"
	"net/http"

	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ErrDuplicateUsername = errors.New("username already exists")

type RegisterHandler struct {
	db              *gorm.DB
	registerService services.RegisterService
}

func NewRegisterHandler(db *gorm.DB, registerService services.RegisterService) *RegisterHandler {
	return &RegisterHandler{db: db, registerService: registerService}
}

// Registration handles POST /api/v1/auth/register.
// TODO(1): Bind a request struct (username, email, password) with validation; 400 on invalid input.
// Do not bind straight into models.User, so clients can't set ID, role or timestamps.
// TODO(1): Hash the password with bcrypt (golang.org/x/crypto/bcrypt) in the service; never store plaintext.
// TODO(1): Generate the user's UUID, save with GORM; return 409 for a duplicate username/email (ErrDuplicateUsername).
// TODO(1): Never return the password hash in the response.
// TODO(5): The service must also assign the default "user" role (insert into user_roles).
func (h *RegisterHandler) Registration(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
