package handlers

import (
	"net/http"

	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RefreshHandler struct {
	db          *gorm.DB
	authService services.AuthService
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewRefreshHandler(db *gorm.DB, authService services.AuthService) *RefreshHandler {
	return &RefreshHandler{db: db, authService: authService}
}

// Refresh handles POST /api/v1/auth/refresh.
// TODO(3): Bind RefreshRequest; find the refresh_token in the tokens table; 401 if unknown or expired.
// TODO(3): Issue a new access_token AND a new refresh_token, with expiries (rubric).
// TODO(3): Rotate: delete or invalidate the used refresh_token so it can't be replayed.
// TODO(5): Re-load the user's current roles/permissions into the new access_token.
func (h *RefreshHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
		"expires_in":    3600,
	})
}
