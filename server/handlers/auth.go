package handlers

import (
	"database/sql"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResp struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "请输入用户名和密�?})
	}
	var userID int
	var hash, role string
	err := h.DB.QueryRow("SELECT id, password_hash, role FROM users WHERE username = ?", req.Username).Scan(&userID, &hash, &role)
	if err == sql.ErrNoRows {
		return c.Status(401).JSON(fiber.Map{"error": "用户名或密码错误"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "服务器错�?})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "用户名或密码错误"})
	}
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(720 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "生成Token失败"})
	}
	return c.JSON(loginResp{
		Token:    tokenStr,
		UserID:   userID,
		Username: req.Username,
		Role:     role,
	})
}

type changePWReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req changePWReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "新密码不能为�?})
	}
	userID := c.Locals("user_id").(int)
	var hash string
	if err := h.DB.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&hash); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "服务器错�?})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.OldPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "原密码错�?})
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "加密失败"})
	}
	if _, err := h.DB.Exec("UPDATE users SET password_hash = ?, updated_at = datetime('now') WHERE id = ?", string(newHash), userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "修改失败"})
	}
	return c.JSON(fiber.Map{"message": "修改成功"})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	var username, role string
	err := h.DB.QueryRow("SELECT username, role FROM users WHERE id = ?", userID).Scan(&username, &role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "服务器错�?})
	}
	return c.JSON(fiber.Map{"user_id": userID, "username": username, "role": role})
}
