package middleware

import (
	"strings"

	"expensio-backend/internal/config"
	"expensio-backend/pkg/cache"
	jwtUtil "expensio-backend/pkg/jwt"
	"expensio-backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header is required")
		}

		tokenString, err := jwtUtil.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		// Validate token
		claims, err := jwtUtil.ValidateToken(tokenString, cfg)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Check if token type is access token
		if claims.TokenType != "access" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token type")
		}

		// Check if token is blacklisted (logged out)
		isBlacklisted, _ := cache.Exists("blacklist:token:" + claims.ID)
		if isBlacklisted {
			return response.Error(c, fiber.StatusUnauthorized, "Token has been revoked")
		}

		// Check if session exists in Redis
		sessionKey := "session:" + claims.UserID
		exists, _ := cache.Exists(sessionKey)
		if !exists {
			return response.Error(c, fiber.StatusUnauthorized, "Session expired or invalid")
		}

		// Set user context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("companyID", claims.CompanyID)
		c.Locals("tokenID", claims.ID)

		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return response.Error(c, fiber.StatusUnauthorized, "User role not found")
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(role, allowedRole) {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "Insufficient permissions")
	}
}

// OptionalAuthMiddleware validates token if present, but doesn't require it
func OptionalAuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenString, err := jwtUtil.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Next()
		}

		claims, err := jwtUtil.ValidateToken(tokenString, cfg)
		if err != nil {
			return c.Next()
		}

		// Set user context if token is valid
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("companyID", claims.CompanyID)

		return c.Next()
	}
}
