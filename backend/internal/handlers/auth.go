package handlers

import (
	"net/http"
	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db          *gorm.DB
	authService services.AuthService
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(db *gorm.DB, authService services.AuthService) *AuthHandler {
	return &AuthHandler{db: db, authService: authService}
}

// Token handles POST /api/v1/auth/login.
// TODO(2): Bind AuthRequest; look the user up by username; verify with services.VerifyPassword.
// TODO(2): Return the same 401 message for "no such user" and "wrong password" (don't leak which one failed).
// TODO(2): Issue a JWT access_token (valid 1 hour) and a refresh_token (valid 1 hour, rubric).
// TODO(2): Store the refresh_token and its expiry in the tokens table (models.Token).
// TODO(2): Respond with access_token, refresh_token and expires_in (seconds).
// TODO(5): Load the user's roles and permissions and embed them in the access_token claims (utils.Claims).
func (h *AuthHandler) Token(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
		"expires_in":    3600,
	})
}
