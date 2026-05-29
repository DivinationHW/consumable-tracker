package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type ProblemTypeHandler struct {
	DB *sql.DB
}

func NewProblemTypeHandler(db *sql.DB) *ProblemTypeHandler {
	return &ProblemTypeHandler{DB: db}
}

func (h *ProblemTypeHandler) List(c *fiber.Ctx) error {
	deviceType := c.Query("device_type")
	var rows *sql.Rows
	var err error

	if deviceType != "" {
		rows, err = h.DB.Query(
			"SELECT id, device_type, name, sort_order, is_default, created_at FROM problem_types WHERE device_type = $1 ORDER BY sort_order, id",
			deviceType,
		)
	} else {
		rows, err = h.DB.Query("SELECT id, device_type, name, sort_order, is_default, created_at FROM problem_types ORDER BY device_type, sort_order, id")
	}
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query problem types"})
	}
	defer rows.Close()

	var items []models.ProblemType
	for rows.Next() {
		var item models.ProblemType
		if err := rows.Scan(&item.ID, &item.DeviceType, &item.Name, &item.SortOrder, &item.IsDefault, &item.CreatedAt); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan problem type"})
		}
		items = append(items, item)
	}
	return c.JSON(items)
}

func (h *ProblemTypeHandler) Create(c *fiber.Ctx) error {
	var req models.ProblemTypeCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Name == "" || req.DeviceType == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "name and device_type are required"})
	}
	if len(req.Name) > 50 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "name too long (max 50)"})
	}

	var item models.ProblemType
	err := h.DB.QueryRow(
		"INSERT INTO problem_types (device_type, name, sort_order) VALUES ($1, $2, $3) RETURNING id, device_type, name, sort_order, is_default, created_at",
		req.DeviceType, req.Name, req.SortOrder,
	).Scan(&item.ID, &item.DeviceType, &item.Name, &item.SortOrder, &item.IsDefault, &item.CreatedAt)
	if err != nil {
		return c.Status(409).JSON(models.ErrorResponse{Error: "problem type already exists for this device type"})
	}
	return c.Status(201).JSON(item)
}

func (h *ProblemTypeHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.ProblemTypeCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	_, err := h.DB.Exec(
		"UPDATE problem_types SET name = $1, sort_order = $2 WHERE id = $3",
		req.Name, req.SortOrder, id,
	)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update problem type"})
	}
	return c.JSON(fiber.Map{"message": "problem type updated"})
}

func (h *ProblemTypeHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var isDefault bool
	h.DB.QueryRow("SELECT is_default FROM problem_types WHERE id = $1", id).Scan(&isDefault)
	if isDefault {
		return c.Status(400).JSON(models.ErrorResponse{Error: "cannot delete default problem type"})
	}

	_, err := h.DB.Exec("DELETE FROM problem_types WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete problem type"})
	}
	return c.JSON(fiber.Map{"message": "problem type deleted"})
}

func (h *ProblemTypeHandler) Sort(c *fiber.Ctx) error {
	var items []models.ProblemTypeCreate
	if err := c.BodyParser(&items); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	for i, item := range items {
		_, err := h.DB.Exec("UPDATE problem_types SET sort_order = $1 WHERE device_type = $2 AND name = $3", i, item.DeviceType, item.Name)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to reorder"})
		}
	}
	return c.JSON(fiber.Map{"message": "reordered successfully"})
}
