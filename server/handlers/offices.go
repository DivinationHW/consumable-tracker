package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type OfficeHandler struct {
	DB *sql.DB
}

func NewOfficeHandler(db *sql.DB) *OfficeHandler {
	return &OfficeHandler{DB: db}
}

func (h *OfficeHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, room_number, device_type, device_model, created_at FROM offices ORDER BY room_number")
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query offices"})
	}
	defer rows.Close()

	var items []models.Office
	for rows.Next() {
		var item models.Office
		if err := rows.Scan(&item.ID, &item.RoomNumber, &item.DeviceType, &item.DeviceModel, &item.CreatedAt); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan office"})
		}
		items = append(items, item)
	}
	return c.JSON(items)
}

func (h *OfficeHandler) Create(c *fiber.Ctx) error {
	var req models.OfficeCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.RoomNumber == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "room_number is required"})
	}
	if len(req.RoomNumber) > 20 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "room_number too long (max 20)"})
	}

	var item models.Office
	err := h.DB.QueryRow(
		"INSERT INTO offices (room_number, device_type, device_model) VALUES ($1, $2, $3) RETURNING id, room_number, device_type, device_model, created_at",
		req.RoomNumber, req.DeviceType, req.DeviceModel,
	).Scan(&item.ID, &item.RoomNumber, &item.DeviceType, &item.DeviceModel, &item.CreatedAt)
	if err != nil {
		return c.Status(409).JSON(models.ErrorResponse{Error: "room number already exists"})
	}
	return c.Status(201).JSON(item)
}

func (h *OfficeHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	var req models.OfficeUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.RoomNumber != nil {
		_, err = h.DB.Exec("UPDATE offices SET room_number = $1 WHERE id = $2", *req.RoomNumber, id)
	} else if req.DeviceType != nil || req.DeviceModel != nil {
		_, err = h.DB.Exec("UPDATE offices SET device_type = COALESCE($1, device_type), device_model = COALESCE($2, device_model) WHERE id = $3",
			req.DeviceType, req.DeviceModel, id)
	}
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update office"})
	}
	return c.JSON(fiber.Map{"message": "office updated"})
}

func (h *OfficeHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	_, err = h.DB.Exec("DELETE FROM offices WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete office"})
	}
	return c.JSON(fiber.Map{"message": "office deleted"})
}
