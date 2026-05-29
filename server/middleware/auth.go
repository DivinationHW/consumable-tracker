package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{"error": "未登�?})
		}
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "登录已过�?})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "无效Token"})
		}
		userID := int(claims["sub"].(float64))
		role := claims["role"].(string)
		c.Locals("user_id", userID)
		c.Locals("role", role)
		return c.Next()
	}
}

func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			return c.Next()
		}
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Next()
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			c.Locals("user_id", int(claims["sub"].(float64)))
			c.Locals("role", claims["role"].(string))
		}
		return c.Next()
	}
}

func RoleAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "无权�?})
		}
		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}
	return c.Cookies("token")
}
