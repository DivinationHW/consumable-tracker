package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type ConsumableHandler struct {
	DB *sql.DB
}

func NewConsumableHandler(db *sql.DB) *ConsumableHandler {
	return &ConsumableHandler{DB: db}
}

func (h *ConsumableHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, name, unit, is_default, created_at FROM consumable_models ORDER BY is_default DESC, name")
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query consumables"})
	}
	defer rows.Close()

	var items []models.Consumable
	for rows.Next() {
		var item models.Consumable
		if err := rows.Scan(&item.ID, &item.Name, &item.Unit, &item.IsDefault, &item.CreatedAt); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan consumable"})
		}
		items = append(items, item)
	}
	return c.JSON(items)
}

func (h *ConsumableHandler) Create(c *fiber.Ctx) error {
	var req models.ConsumableCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Name == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "name is required"})
	}
	if len(req.Name) > 100 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "name too long (max 100)"})
	}
	if req.Unit == "" {
		req.Unit = "个"
	}

	var item models.Consumable
	err := h.DB.QueryRow(
		"INSERT INTO consumable_models (name, unit) VALUES ($1, $2) RETURNING id, name, unit, is_default, created_at",
		req.Name, req.Unit,
	).Scan(&item.ID, &item.Name, &item.Unit, &item.IsDefault, &item.CreatedAt)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create consumable"})
	}
	return c.Status(201).JSON(item)
}

func (h *ConsumableHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	var req models.ConsumableCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	_, err = h.DB.Exec("UPDATE consumable_models SET name = $1, unit = $2 WHERE id = $3", req.Name, req.Unit, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update consumable"})
	}
	return c.JSON(fiber.Map{"message": "consumable updated"})
}

func (h *ConsumableHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	var isDefault bool
	h.DB.QueryRow("SELECT is_default FROM consumable_models WHERE id = $1", id).Scan(&isDefault)
	if isDefault {
		return c.Status(400).JSON(models.ErrorResponse{Error: "cannot delete default consumable"})
	}

	_, err = h.DB.Exec("DELETE FROM consumable_models WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete consumable"})
	}
	return c.JSON(fiber.Map{"message": "consumable deleted"})
}
