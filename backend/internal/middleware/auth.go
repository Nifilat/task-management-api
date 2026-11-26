package middleware

import (
	"net/http"
	"strings"
	"task-manager/backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the access_token on each request and stores its claims on the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// RequirePermission allows the request only if the token carries the given permission (RBAC + resource-level ABAC).
// TODO(6): Read claims from the context; respond 403 if the permission resource:action is missing.
// TODO(6): Ownership rule (ABAC): a "user" may only touch their own profile/tasks (claims.user_id == owner);
// an admin may act on anyone's. Decide where this check lives (middleware vs handler/service).
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // pass-through until implemented
	}
}

// RequireRole allows the request only for the given role, e.g. RequireRole("admin") for /users management.
// TODO(6): Respond 403 when the role is not in the token's roles claim.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // pass-through until implemented
	}
}
