package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "unauthorized: missing user_id",
			})
		}

		role := c.Get("X-Role")
		if role == "" {
			role = "user"
		}

		c.Locals("userID", userID)
		c.Locals("role", role)

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("userID").(string); ok {
		return userID
	}
	return ""
}

func GetRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("role").(string); ok {
		return role
	}
	return "user"
}

func IsAdmin(c *fiber.Ctx) bool {
	return GetRole(c) == "admin"
}

func GetTargetUserID(c *fiber.Ctx) string {
	currentUserID := GetUserID(c)

	if IsAdmin(c) {
		targetUserID := c.FormValue("target_user_id")
		if targetUserID != "" {
			return targetUserID
		}
	}

	return currentUserID
}
