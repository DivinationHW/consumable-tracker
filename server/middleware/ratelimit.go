package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	lm := &ipLimiter{
		visitors: make(map[string]*rate.Limiter),
	}
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			lm.mu.Lock()
			for k, v := range lm.visitors {
				if v.Tokens() >= float64(b) {
					delete(lm.visitors, k)
				}
			}
			lm.mu.Unlock()
		}
	}()
	return lm
}

func (l *ipLimiter) get(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter, exists := l.visitors[key]
	if !exists {
		limiter = rate.NewLimiter(1, 5)
		l.visitors[key] = limiter
	}
	return limiter
}

var loginLimiter = newIPLimiter(1, 5)

func LoginRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		ip := c.IP()
		if username == "" {
			username = c.Get("X-Forwarded-For", "unknown")
		}
		key := username + ":" + ip
		if !loginLimiter.get(key).Allow() {
			return c.Status(429).JSON(fiber.Map{"error": "登录尝试过多，请5分钟后再�?})
		}
		return c.Next()
	}
}

var publicLimiter = newIPLimiter(rate.Limit(10.0/3600.0), 10)

func PublicRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !publicLimiter.get(c.IP()).Allow() {
			return c.Status(429).JSON(fiber.Map{"error": "请求过于频繁"})
		}
		return c.Next()
	}
}
