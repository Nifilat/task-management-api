package middleware

import "github.com/gin-gonic/gin"

// AuthMiddleware verifies the access_token on each request and stores its claims on the context.
// TODO(6): Read the "Authorization: Bearer <token>" header; respond 401 if missing or malformed.
// TODO(6): Validate the signature, issuer and expiry with utils.ParseJWT; respond 401 if invalid or expired.
// TODO(6): c.Set("claims", claims) (user_id, roles, permissions) for handlers and later middleware, then c.Next().
// TODO(6): Wire it in main.go with .Use() on the protected route groups (not on /auth/register, /login, /refresh).
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // pass-through until implemented: every route is OPEN
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
