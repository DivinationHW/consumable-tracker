package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, username, role, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	type userResp struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	users := []userResp{}
	for rows.Next() {
		var u userResp
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return c.JSON(users)
}

type createUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req createUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "用户名和密码不能为空"})
	}
	if req.Role != "admin" && req.Role != "readonly" {
		req.Role = "readonly"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "加密失败"})
	}
	result, err := h.DB.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)", req.Username, string(hash), req.Role)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "用户名已存在"})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id, "message": "创建成功"})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	userID := c.Locals("user_id").(int)
	if id == userID {
		return c.Status(400).JSON(fiber.Map{"error": "不能删除自己"})
	}
	result, err := h.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "用户不存�?})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

type updateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var req updateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Username != "" {
		if _, err := h.DB.Exec("UPDATE users SET username = ?, updated_at = datetime('now') WHERE id = ?", req.Username, id); err != nil {
			return c.Status(409).JSON(fiber.Map{"error": "用户名已存在"})
		}
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "加密失败"})
		}
		h.DB.Exec("UPDATE users SET password_hash = ?, updated_at = datetime('now') WHERE id = ?", string(hash), id)
	}
	if req.Role != "" {
		h.DB.Exec("UPDATE users SET role = ?, updated_at = datetime('now') WHERE id = ?", req.Role, id)
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}
