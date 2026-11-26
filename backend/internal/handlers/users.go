package handlers

import (
	"net/http"
	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"

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

// GetUserProfile handles GET /api/v1/users/profile (the caller's own profile).
// TODO(8): SECURITY: services.GetUserProfileMalicious builds SQL with fmt.Sprintf and is injectable.
// Use the parameterized GetUserProfile instead, and delete the Malicious method (rubric: SQL injection fix).
// TODO(8): Take the user id from the token claims, not from the request. Never return the password hash.
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	userID, err := uuid.FromString(value.(*utils.Claims).UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user, err := h.userService.GetUserProfile(h.db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load profile"})
		return
	}
	if user.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserProfileByUserId handles GET /api/v1/users/profile/:user_id.
// TODO(8): Validate :user_id as a UUID, use the parameterized query, 404 if not found.
// TODO(8): ABAC: a user may read only their own profile; an admin may read anyone's.
func (h *UserHandler) GetUserProfileByUserId(c *gin.Context) {
	c.JSON(http.StatusOK, "user profile by user ID")
}

// GetUsers handles GET /api/v1/users (admin only: user:read).
// TODO(8): Restrict to admins; never include password hashes in the response.
func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, "users list")
}

// DeleteUser handles DELETE /api/v1/users/:user_id (admin only: user:delete).
// TODO(8): Validate the UUID, 404 if missing. Decide whether an admin may delete their own account.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
