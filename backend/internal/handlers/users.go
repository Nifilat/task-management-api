package handlers

import (
	"errors"
	"net/http"
	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db          *gorm.DB
	userService services.UserService
}

func NewUserHandler(db *gorm.DB, userService services.UserService) *UserHandler {
	return &UserHandler{db: db, userService: userService}
}

// GetUserProfile handles GET /api/v1/users/profile (the caller's own profile; the id comes from the token).
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	_, userID, ok := callerFromContext(c)
	if !ok {
		return
	}

	user, err := h.userService.GetUserProfile(h.db, userID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserProfileByUserId handles GET /api/v1/users/profile/:user_id.
// The route's RequireOwnerOrAdmin middleware ensures the caller is that user or an admin.
func (h *UserHandler) GetUserProfileByUserId(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userService.GetUserProfile(h.db, userID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers handles GET /api/v1/users (admin only: the route requires user:read).
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers(h.db)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteUser handles DELETE /api/v1/users/:user_id (admin only: the route requires user:delete).
// An admin cannot delete their own account, so the platform is never left without one.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	_, callerID, ok := callerFromContext(c)
	if !ok {
		return
	}
	if userID == callerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot delete your own account"})
		return
	}

	if err := h.userService.DeleteUser(h.db, userID); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func handleUserError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user request"})
	}
}
