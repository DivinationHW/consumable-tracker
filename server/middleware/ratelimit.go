package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	maxAttempts int
	window     time.Duration
}

func newRateLimiter(maxAttempts int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		attempts:    make(map[string][]time.Time),
		maxAttempts: maxAttempts,
		window:      window,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.attempts {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.attempts, key)
			} else {
				rl.attempts[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.attempts[key]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxAttempts {
		rl.attempts[key] = valid
		return false
	}

	rl.attempts[key] = append(valid, now)
	return true
}

var loginLimiter = newRateLimiter(5, 5*time.Minute)
var publicLimiter = newRateLimiter(10, 1*time.Hour)
var qrcodeLimiter = newRateLimiter(10, 1*time.Minute)

func LoginRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		if username == "" {
			username = c.Query("username")
		}
		clientIP := c.IP()
		key := username + ":" + clientIP

		if !loginLimiter.allow(key) {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many login attempts, try again in 5 minutes",
			})
		}
		return c.Next()
	}
}

func PublicRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !publicLimiter.allow(c.IP()) {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests, try again later",
			})
		}
		return c.Next()
	}
}

func QRCodeRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !qrcodeLimiter.allow(c.IP()) {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many qrcode generation requests, try again later",
			})
		}
		return c.Next()
	}
}
