package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}

func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}

func AdminOrReadonly() fiber.Handler {
	return RequireRole("admin", "readonly")
}
