package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type NoteHandler struct {
	DB *sql.DB
}

func (h *NoteHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.Note{}
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		list = append(list, n)
	}
	return c.JSON(list)
}

func (h *NoteHandler) Create(c *fiber.Ctx) error {
	var n models.Note
	if err := c.BodyParser(&n); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if n.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "标题不能为空"})
	}
	result, err := h.DB.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", n.Title, n.Content)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建失败"})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *NoteHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var n models.Note
	if err := c.BodyParser(&n); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	_, err = h.DB.Exec("UPDATE notes SET title=?, content=?, updated_at=datetime('now') WHERE id=?", n.Title, n.Content, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}

func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	_, err = h.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}
