package middleware

import (
	"net/http"
	"slices"
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

// getClaims returns the claims stored by AuthMiddleware; it replies 401 and returns false when they are missing.
func getClaims(c *gin.Context) (*utils.Claims, bool) {
	value, exists := c.Get("claims")
	claims, ok := value.(*utils.Claims)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		c.Abort()
		return nil, false
	}
	return claims, true
}

// RequirePermission allows the request only if the token carries the given permission (RBAC).
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := getClaims(c)
		if !ok {
			return
		}

		for _, p := range claims.Permissions {
			if p.Resource == resource && slices.Contains(p.Actions, action) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "Insufficient permissions",
			"required": resource + ":" + action,
		})
		c.Abort()
	}
}

// RequireRole allows the request only for the given role, e.g. RequireRole("admin").
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := getClaims(c)
		if !ok {
			return
		}

		if !slices.Contains(claims.Roles, role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":     "Insufficient permissions",
				"your_role": claims.Roles,
				"required":  role,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnerOrAdmin is the ownership check (ABAC): the user id in the URL must be the caller's own, unless the caller is an admin.
func RequireOwnerOrAdmin(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := getClaims(c)
		if !ok {
			return
		}

		if claims.IsAdmin || strings.EqualFold(claims.UserID, c.Param(param)) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own data"})
		c.Abort()
	}
}
