package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/server/models"
)

type ConsumableHandler struct {
	DB *sql.DB
}

func (h *ConsumableHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, name, unit, is_default, created_at FROM consumable_models ORDER BY is_default DESC, name")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.ConsumableModel{}
	for rows.Next() {
		var m models.ConsumableModel
		if err := rows.Scan(&m.ID, &m.Name, &m.Unit, &m.IsDefault, &m.CreatedAt); err != nil {
			continue
		}
		list = append(list, m)
	}
	return c.JSON(list)
}

func (h *ConsumableHandler) Create(c *fiber.Ctx) error {
	var m models.ConsumableModel
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if m.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "名称不能为空"})
	}
	if m.Unit == "" {
		m.Unit = "�?
	}
	result, err := h.DB.Exec("INSERT INTO consumable_models (name, unit) VALUES (?, ?)", m.Name, m.Unit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建失败"})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *ConsumableHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var m models.ConsumableModel
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if m.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "名称不能为空"})
	}
	if m.Unit == "" {
		m.Unit = "�?
	}
	_, err = h.DB.Exec("UPDATE consumable_models SET name = ?, unit = ? WHERE id = ?", m.Name, m.Unit, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}

func (h *ConsumableHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	result, err := h.DB.Exec("DELETE FROM consumable_models WHERE id = ? AND is_default = 0", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "默认耗材无法删除或不存在"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}
