package handlers

import (
	"net/http"
	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, "user profile")
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
