package handlers

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type RecordHandler struct {
	DB *sql.DB
}

func NewRecordHandler(db *sql.DB) *RecordHandler {
	return &RecordHandler{DB: db}
}

func (h *RecordHandler) List(c *fiber.Ctx) error {
	query := `SELECT r.id, r.office_id, r.consumable_id, r.quantity, 
		to_char(r.usage_date, 'YYYY-MM-DD'), COALESCE(r.note, ''), r.created_at, r.updated_at,
		o.room_number, c.name
		FROM usage_records r
		JOIN offices o ON r.office_id = o.id
		JOIN consumable_models c ON r.consumable_id = c.id
		WHERE 1=1`

	var args []interface{}
	argIdx := 1

	if officeID := c.Query("office"); officeID != "" {
		query += fmt.Sprintf(" AND r.office_id = $%d", argIdx)
		args = append(args, officeID)
		argIdx++
	}
	if consumableID := c.Query("consumable"); consumableID != "" {
		query += fmt.Sprintf(" AND r.consumable_id = $%d", argIdx)
		args = append(args, consumableID)
		argIdx++
	}
	if startDate := c.Query("start"); startDate != "" {
		query += fmt.Sprintf(" AND r.usage_date >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate := c.Query("end"); endDate != "" {
		query += fmt.Sprintf(" AND r.usage_date <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}

	query += " ORDER BY r.usage_date DESC, r.created_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query records"})
	}
	defer rows.Close()

	var records []models.UsageRecord
	for rows.Next() {
		var r models.UsageRecord
		if err := rows.Scan(&r.ID, &r.OfficeID, &r.ConsumableID, &r.Quantity,
			&r.UsageDate, &r.Note, &r.CreatedAt, &r.UpdatedAt,
			&r.OfficeNumber, &r.ConsumableName); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan record"})
		}
		records = append(records, r)
	}
	return c.JSON(records)
}

func (h *RecordHandler) Create(c *fiber.Ctx) error {
	var req models.UsageRecordCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Quantity <= 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "quantity must be positive"})
	}
	if len(req.Note) > 200 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "note too long (max 200)"})
	}

	var r models.UsageRecord
	err := h.DB.QueryRow(
		`INSERT INTO usage_records (office_id, consumable_id, quantity, usage_date, note)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, $4, created_at, updated_at,
		(SELECT room_number FROM offices WHERE id = $1),
		(SELECT name FROM consumable_models WHERE id = $2)`,
		req.OfficeID, req.ConsumableID, req.Quantity, req.UsageDate, req.Note,
	).Scan(&r.ID, &r.UsageDate, &r.CreatedAt, &r.UpdatedAt, &r.OfficeNumber, &r.ConsumableName)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create record"})
	}
	r.OfficeID = req.OfficeID
	r.ConsumableID = req.ConsumableID
	r.Quantity = req.Quantity
	r.Note = req.Note

	return c.Status(201).JSON(r)
}

func (h *RecordHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.UsageRecordCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	_, err := h.DB.Exec(
		`UPDATE usage_records SET office_id = $1, consumable_id = $2, quantity = $3, usage_date = $4, note = $5, updated_at = NOW()
		WHERE id = $6`,
		req.OfficeID, req.ConsumableID, req.Quantity, req.UsageDate, req.Note, id,
	)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update record"})
	}
	return c.JSON(fiber.Map{"message": "record updated"})
}

func (h *RecordHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM usage_records WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete record"})
	}
	return c.JSON(fiber.Map{"message": "record deleted"})
}
