package handlers

import (
	"errors"
	"net/http"

	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"

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
func (h *RefreshHandler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	accessToken, refreshToken, err := h.authService.RotateRefreshToken(h.db, req.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    int(utils.AccessTokenTTL.Seconds()),
	})
}
