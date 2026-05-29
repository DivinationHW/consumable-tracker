package middleware

import "github.com/gofiber/fiber/v2"

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")
		return c.Next()
	}
}

func HSTS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		return c.Next()
	}
}
