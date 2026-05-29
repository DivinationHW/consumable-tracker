package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/server/models"
)

type ProblemTypeHandler struct {
	DB *sql.DB
}

func (h *ProblemTypeHandler) List(c *fiber.Ctx) error {
	deviceType := c.Query("device_type")
	query := "SELECT id, device_type, name, sort_order, is_default, created_at FROM problem_types"
	args := []interface{}{}
	if deviceType != "" {
		query += " WHERE device_type = ?"
		args = append(args, deviceType)
	}
	query += " ORDER BY device_type, sort_order, name"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.ProblemType{}
	for rows.Next() {
		var p models.ProblemType
		if err := rows.Scan(&p.ID, &p.DeviceType, &p.Name, &p.SortOrder, &p.IsDefault, &p.CreatedAt); err != nil {
			continue
		}
		list = append(list, p)
	}
	return c.JSON(list)
}

func (h *ProblemTypeHandler) Create(c *fiber.Ctx) error {
	var p models.ProblemType
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if p.DeviceType == "" || p.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "参数不完�?})
	}
	result, err := h.DB.Exec("INSERT INTO problem_types (device_type, name, sort_order) VALUES (?, ?, ?)", p.DeviceType, p.Name, p.SortOrder)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "该设备类型下已存在同名故障类�?})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *ProblemTypeHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var p models.ProblemType
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	_, err = h.DB.Exec("UPDATE problem_types SET device_type=?, name=?, sort_order=? WHERE id=?", p.DeviceType, p.Name, p.SortOrder, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}

func (h *ProblemTypeHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	result, err := h.DB.Exec("DELETE FROM problem_types WHERE id = ? AND is_default = 0", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "默认故障类型无法删除或不存在"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

func (h *ProblemTypeHandler) GetDeviceTypes(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT DISTINCT device_type FROM problem_types ORDER BY device_type")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	types := []string{}
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			continue
		}
		types = append(types, t)
	}
	return c.JSON(types)
}
