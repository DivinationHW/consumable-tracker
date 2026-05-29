package handlers

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/middleware"
	"consumable-tracker/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "username and password are required"})
	}

	var user models.User
	err := h.DB.QueryRow(
		"SELECT id, username, password_hash, role FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(models.ErrorResponse{Error: "invalid username or password"})
	}
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		hashed := sha256.Sum256([]byte(req.Password))
		if fmt.Sprintf("%x", hashed) != user.PasswordHash {
			return c.Status(401).JSON(models.ErrorResponse{Error: "invalid username or password"})
		}
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to generate token"})
	}

	return c.JSON(models.TokenResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "old and new passwords are required"})
	}
	if len(req.NewPassword) < 8 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "password must be at least 8 characters"})
	}

	var passwordHash string
	err := h.DB.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.OldPassword)); err != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: "invalid old password"})
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to hash password"})
	}

	_, err = h.DB.Exec("UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2", string(newHash), userID)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update password"})
	}

	return c.JSON(fiber.Map{"message": "password updated successfully"})
}
