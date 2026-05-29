package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"consumable-tracker/models"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, username, role, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query users"})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan user"})
		}
		users = append(users, u)
	}
	return c.JSON(users)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.UserCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "username and password are required"})
	}
	if len(req.Password) < 8 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "password must be at least 8 characters"})
	}
	if req.Role == "" {
		req.Role = "readonly"
	}
	if req.Role != "admin" && req.Role != "readonly" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "role must be admin or readonly"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to hash password"})
	}

	var user models.User
	err = h.DB.QueryRow(
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id, username, role, created_at",
		req.Username, string(hash), req.Role,
	).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		return c.Status(409).JSON(models.ErrorResponse{Error: "username already exists"})
	}

	return c.Status(201).JSON(user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	var req models.UserUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.Username != nil {
		_, err = h.DB.Exec("UPDATE users SET username = $1, updated_at = NOW() WHERE id = $2", *req.Username, id)
		if err != nil {
			return c.Status(409).JSON(models.ErrorResponse{Error: "username already exists"})
		}
	}

	if req.Password != nil {
		if len(*req.Password) < 8 {
			return c.Status(400).JSON(models.ErrorResponse{Error: "password must be at least 8 characters"})
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to hash password"})
		}
		_, err = h.DB.Exec("UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2", string(hash), id)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update password"})
		}
	}

	return c.JSON(fiber.Map{"message": "user updated successfully"})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	currentUserID := c.Locals("user_id").(int)
	if id == currentUserID {
		return c.Status(400).JSON(models.ErrorResponse{Error: "cannot delete yourself"})
	}

	_, err = h.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete user"})
	}
	return c.JSON(fiber.Map{"message": "user deleted successfully"})
}
