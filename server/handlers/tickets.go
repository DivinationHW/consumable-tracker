package handlers

import (
	"database/sql"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type TicketHandler struct {
	DB       *sql.DB
	notifier *WebSocketHub
}

func NewTicketHandler(db *sql.DB, notifier *WebSocketHub) *TicketHandler {
	return &TicketHandler{DB: db, notifier: notifier}
}

func (h *TicketHandler) Submit(c *fiber.Ctx) error {
	var req models.TicketCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.ProblemType == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "problem_type is required"})
	}
	if len(req.Description) > 500 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "description too long (max 500)"})
	}
	if len(req.Contact) > 100 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "contact too long (max 100)"})
	}

	var officeExists bool
	h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM offices WHERE id = $1)", req.OfficeID).Scan(&officeExists)
	if !officeExists {
		return c.Status(400).JSON(models.ErrorResponse{Error: "office not found"})
	}

	var ticket models.Ticket
	err := h.DB.QueryRow(`INSERT INTO tickets (office_id, problem_type, description, contact, device_type, device_model, status)
		VALUES ($1, $2, $3, $4,
			COALESCE((SELECT device_type FROM offices WHERE id = $1), ''),
			COALESCE((SELECT device_model FROM offices WHERE id = $1), ''),
			'pending')
		RETURNING id, office_id, device_type, device_model, problem_type, COALESCE(description,''), COALESCE(contact,''),
			status, created_at, updated_at,
			(SELECT room_number FROM offices WHERE id = $1)`,
		req.OfficeID, req.ProblemType, req.Description, req.Contact,
	).Scan(&ticket.ID, &ticket.OfficeID, &ticket.DeviceType, &ticket.DeviceModel,
		&ticket.ProblemType, &ticket.Description, &ticket.Contact,
		&ticket.Status, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.RoomNumber)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create ticket"})
	}

	if h.notifier != nil {
		data, _ := json.Marshal(fiber.Map{"type": "new_ticket", "ticket": ticket})
		h.notifier.Broadcast("admin", string(data))
	}

	return c.Status(201).JSON(ticket)
}

func (h *TicketHandler) GetStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var ticket models.Ticket
	err := h.DB.QueryRow(`SELECT id, problem_type, status, created_at, 
		(SELECT room_number FROM offices WHERE id = tickets.office_id)
		FROM tickets WHERE id = $1`, id,
	).Scan(&ticket.ID, &ticket.ProblemType, &ticket.Status, &ticket.CreatedAt, &ticket.RoomNumber)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ErrorResponse{Error: "ticket not found"})
	}
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query ticket"})
	}
	return c.JSON(ticket)
}

func (h *TicketHandler) OfficeProblemTypes(c *fiber.Ctx) error {
	officeID := c.Params("office_id")
	var deviceType string
	h.DB.QueryRow("SELECT COALESCE(device_type, 'other') FROM offices WHERE id = $1", officeID).Scan(&deviceType)

	rows, err := h.DB.Query("SELECT id, name, is_default FROM problem_types WHERE device_type = $1 ORDER BY sort_order, id", deviceType)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query problem types"})
	}
	defer rows.Close()

	type ProblemTypeOption struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
	}

	var options []ProblemTypeOption
	for rows.Next() {
		var opt ProblemTypeOption
		if rows.Scan(&opt.ID, &opt.Name, &opt.IsDefault) == nil {
			options = append(options, opt)
		}
	}
	return c.JSON(fiber.Map{"office_id": officeID, "device_type": deviceType, "problem_types": options})
}

func (h *TicketHandler) List(c *fiber.Ctx) error {
	query := `SELECT id, office_id, device_type, device_model, problem_type, 
		COALESCE(description,''), COALESCE(contact,''), status,
		COALESCE(consumable_used,''), COALESCE(consumable_quantity,0),
		handled_by_user_id, COALESCE(handle_note,''), created_at, updated_at,
		(SELECT room_number FROM offices WHERE id = tickets.office_id)
		FROM tickets WHERE 1=1`

	var args []interface{}
	argIdx := 1

	if status := c.Query("status"); status != "" {
		query += " AND status = $1"
		args = append(args, status)
		argIdx++
	}
	if officeID := c.Query("office"); officeID != "" {
		query += " AND office_id = $" + string(rune('0'+argIdx))
		args = append(args, officeID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query tickets"})
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.OfficeID, &t.DeviceType, &t.DeviceModel,
			&t.ProblemType, &t.Description, &t.Contact, &t.Status,
			&t.ConsumableUsed, &t.ConsumableQuantity,
			&t.HandledByUserID, &t.HandleNote, &t.CreatedAt, &t.UpdatedAt, &t.RoomNumber); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan ticket"})
		}
		tickets = append(tickets, t)
	}
	return c.JSON(tickets)
}

func (h *TicketHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.TicketStatusUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Status != "pending" && req.Status != "processing" && req.Status != "completed" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid status"})
	}

	userID := c.Locals("user_id").(int)
	_, err := h.DB.Exec("UPDATE tickets SET status = $1, handled_by_user_id = $2, updated_at = NOW() WHERE id = $3",
		req.Status, userID, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update ticket status"})
	}

	if h.notifier != nil {
		data, _ := json.Marshal(fiber.Map{"type": "ticket_updated", "ticket_id": id, "status": req.Status})
		h.notifier.Broadcast("admin", string(data))
	}

	return c.JSON(fiber.Map{"message": "status updated"})
}

func (h *TicketHandler) Complete(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.TicketComplete
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	userID := c.Locals("user_id").(int)

	tx, err := h.DB.Begin()
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to start transaction"})
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE tickets SET status = 'completed', consumable_used = $1, consumable_quantity = $2, 
		handle_note = $3, handled_by_user_id = $4, updated_at = NOW() WHERE id = $5`,
		req.ConsumableUsed, req.ConsumableQuantity, req.HandleNote, userID, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update ticket"})
	}

	if req.ConsumableUsed != "" && req.ConsumableQuantity > 0 {
		var officeID int
		var consumableID int
		tx.QueryRow("SELECT office_id FROM tickets WHERE id = $1", id).Scan(&officeID)
		tx.QueryRow("SELECT id FROM consumable_models WHERE name ILIKE $1 LIMIT 1", "%"+req.ConsumableUsed+"%").Scan(&consumableID)
		if consumableID == 0 {
			err = tx.QueryRow("INSERT INTO consumable_models (name) VALUES ($1) RETURNING id", req.ConsumableUsed).Scan(&consumableID)
			if err != nil {
				return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create consumable"})
			}
		}
		_, err = tx.Exec("INSERT INTO usage_records (office_id, consumable_id, quantity, usage_date, note) VALUES ($1, $2, $3, CURRENT_DATE, $4)",
			officeID, consumableID, req.ConsumableQuantity, "工单 #"+id[:8]+" 消耗")
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create usage record"})
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to commit"})
	}

	if h.notifier != nil {
		data, _ := json.Marshal(fiber.Map{"type": "ticket_completed", "ticket_id": id})
		h.notifier.Broadcast("admin", string(data))
	}

	return c.JSON(fiber.Map{"message": "ticket completed"})
}

func (h *TicketHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM tickets WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete ticket"})
	}
	return c.JSON(fiber.Map{"message": "ticket deleted"})
}
