package handlers

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"consumable-tracker/server/models"
)

type TicketHandler struct {
	DB        *sql.DB
	Broadcast func(msg interface{})
}

func (h *TicketHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	officeID := c.Query("office_id")
	query := `SELECT t.id, t.office_id, t.device_type, t.device_model, t.problem_type,
		t.description, t.contact, t.status, t.consumable_used, t.consumable_quantity,
		t.handled_by_user_id, t.handle_note, t.created_at, t.updated_at,
		o.room_number, COALESCE(u.username,'')
		FROM tickets t
		JOIN offices o ON t.office_id = o.id
		LEFT JOIN users u ON t.handled_by_user_id = u.id
		WHERE 1=1`
	args := []interface{}{}
	if status != "" {
		query += " AND t.status = ?"
		args = append(args, status)
	}
	if officeID != "" {
		query += " AND t.office_id = ?"
		args = append(args, officeID)
	}
	query += " ORDER BY t.created_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.Ticket{}
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.OfficeID, &t.DeviceType, &t.DeviceModel, &t.ProblemType,
			&t.Description, &t.Contact, &t.Status, &t.ConsumableUsed, &t.ConsumableQty,
			&t.HandledByUserID, &t.HandleNote, &t.CreatedAt, &t.UpdatedAt,
			&t.OfficeName, &t.HandledByUser); err != nil {
			continue
		}
		list = append(list, t)
	}
	return c.JSON(list)
}

type createTicketReq struct {
	OfficeID    int    `json:"office_id"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	ProblemType string `json:"problem_type"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
}

func (h *TicketHandler) Create(c *fiber.Ctx) error {
	var req createTicketReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if qrOfficeID, ok := c.Locals("qr_office_id").(int); ok && qrOfficeID > 0 {
		req.OfficeID = qrOfficeID
	}
	if req.DeviceType == "" {
		if dt, ok := c.Locals("qr_device_type").(string); ok {
			req.DeviceType = dt
		}
	}
	if req.DeviceModel == "" {
		if dm, ok := c.Locals("qr_device_model").(string); ok {
			req.DeviceModel = dm
		}
	}
	if req.OfficeID == 0 || req.ProblemType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "参数不完整"})
	}
	id := uuid.New().String()
	_, err := h.DB.Exec(`INSERT INTO tickets (id, office_id, device_type, device_model, problem_type, description, contact, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')`, id, req.OfficeID, req.DeviceType, req.DeviceModel,
		req.ProblemType, req.Description, req.Contact)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建失败"})
	}
	if h.Broadcast != nil {
		h.Broadcast(fiber.Map{"type": "ticket_new", "id": id})
	}
	return c.Status(201).JSON(fiber.Map{"id": id})
}
	if req.OfficeID == 0 || req.ProblemType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "参数不完�?})
	}
	id := uuid.New().String()
	_, err := h.DB.Exec(`INSERT INTO tickets (id, office_id, device_type, device_model, problem_type, description, contact, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')`, id, req.OfficeID, req.DeviceType, req.DeviceModel,
		req.ProblemType, req.Description, req.Contact)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建失败"})
	}
	if h.Broadcast != nil {
		h.Broadcast(fiber.Map{"type": "ticket_new", "id": id})
	}
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *TicketHandler) CreatePublic(c *fiber.Ctx) error {
	qr := c.Query("qr")
	if qr != "" {
		var officeID int
		var deviceType, deviceModel string
		err := h.DB.QueryRow("SELECT COALESCE(office_id,0), device_type, device_model FROM qr_codes WHERE code = ?", qr).
			Scan(&officeID, &deviceType, &deviceModel)
		if err == nil && officeID > 0 {
			c.Locals("qr_office_id", officeID)
			c.Locals("qr_device_type", deviceType)
			c.Locals("qr_device_model", deviceModel)
		}
	}
	return h.Create(c)
}

func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var t models.Ticket
	err := h.DB.QueryRow(`SELECT t.id, t.office_id, t.device_type, t.device_model, t.problem_type,
		t.description, t.contact, t.status, t.consumable_used, t.consumable_quantity,
		t.handled_by_user_id, t.handle_note, t.created_at, t.updated_at,
		o.room_number, COALESCE(u.username,'')
		FROM tickets t
		JOIN offices o ON t.office_id = o.id
		LEFT JOIN users u ON t.handled_by_user_id = u.id
		WHERE t.id = ?`, id).
		Scan(&t.ID, &t.OfficeID, &t.DeviceType, &t.DeviceModel, &t.ProblemType,
			&t.Description, &t.Contact, &t.Status, &t.ConsumableUsed, &t.ConsumableQty,
			&t.HandledByUserID, &t.HandleNote, &t.CreatedAt, &t.UpdatedAt,
			&t.OfficeName, &t.HandledByUser)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "工单不存�?})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	return c.JSON(t)
}

func (h *TicketHandler) GetPublicByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var status, problemType, createdAt, room, desc, contact string
	err := h.DB.QueryRow(`SELECT t.status, t.problem_type, t.created_at, o.room_number,
		t.description, t.contact
		FROM tickets t JOIN offices o ON t.office_id = o.id WHERE t.id = ?`, id).
		Scan(&status, &problemType, &createdAt, &room, &desc, &contact)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "工单不存�?})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	return c.JSON(fiber.Map{
		"status":       status,
		"problem_type": problemType,
		"created_at":   createdAt,
		"office_name":  room,
		"description":  desc,
		"contact":      contact,
	})
}

type processTicketReq struct {
	Status    string `json:"status"`
	ConsumableUsed string `json:"consumable_used"`
	ConsumableQty  int    `json:"consumable_quantity"`
	HandleNote     string `json:"handle_note"`
}

func (h *TicketHandler) Process(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(int)
	var req processTicketReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Status != "processing" && req.Status != "completed" {
		return c.Status(400).JSON(fiber.Map{"error": "无效状�?})
	}
	var currentStatus string
	err := h.DB.QueryRow("SELECT status FROM tickets WHERE id = ?", id).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "工单不存�?})
	}
	if currentStatus == "completed" {
		return c.Status(400).JSON(fiber.Map{"error": "工单已完成，无法修改"})
	}
	_, err = h.DB.Exec(`UPDATE tickets SET status=?, consumable_used=?, consumable_quantity=?,
		handle_note=?, handled_by_user_id=?, updated_at=datetime('now') WHERE id=?`,
		req.Status, req.ConsumableUsed, req.ConsumableQty, req.HandleNote, userID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	if h.Broadcast != nil {
		h.Broadcast(fiber.Map{"type": "ticket_update", "id": id, "status": req.Status})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}
