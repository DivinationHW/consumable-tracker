package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/server/models"
)

type OfficeHandler struct {
	DB *sql.DB
}

func (h *OfficeHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, room_number, device_type, device_model, created_at FROM offices ORDER BY room_number")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.Office{}
	for rows.Next() {
		var o models.Office
		if err := rows.Scan(&o.ID, &o.RoomNumber, &o.DeviceType, &o.DeviceModel, &o.CreatedAt); err != nil {
			continue
		}
		list = append(list, o)
	}
	return c.JSON(list)
}

func (h *OfficeHandler) Create(c *fiber.Ctx) error {
	var o models.Office
	if err := c.BodyParser(&o); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if o.RoomNumber == "" {
		return c.Status(400).JSON(fiber.Map{"error": "房间号不能为�?})
	}
	result, err := h.DB.Exec("INSERT INTO offices (room_number, device_type, device_model) VALUES (?, ?, ?)", o.RoomNumber, o.DeviceType, o.DeviceModel)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "房间号已存在"})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *OfficeHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var o models.Office
	if err := c.BodyParser(&o); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if o.RoomNumber == "" {
		return c.Status(400).JSON(fiber.Map{"error": "房间号不能为�?})
	}
	_, err = h.DB.Exec("UPDATE offices SET room_number = ?, device_type = ?, device_model = ? WHERE id = ?", o.RoomNumber, o.DeviceType, o.DeviceModel, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}

func (h *OfficeHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	result, err := h.DB.Exec("DELETE FROM offices WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "不存�?})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

func (h *OfficeHandler) GetDeviceModels(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT DISTINCT device_model FROM offices WHERE device_model != '' ORDER BY device_model")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	models := []string{}
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			continue
		}
		models = append(models, m)
	}
	return c.JSON(models)
}
