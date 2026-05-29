package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"image/png"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"consumable-tracker/models"
)

type QRCodeHandler struct {
	DB *sql.DB
}

func NewQRCodeHandler(db *sql.DB) *QRCodeHandler {
	return &QRCodeHandler{DB: db}
}

func (h *QRCodeHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query(`SELECT q.id, q.code, q.office_id, q.device_type, q.device_model, q.created_at,
		COALESCE(o.room_number, ''), q.office_id IS NOT NULL as is_configured
		FROM qr_codes q LEFT JOIN offices o ON q.office_id = o.id
		ORDER BY q.created_at DESC`)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query qrcodes"})
	}
	defer rows.Close()

	var items []models.QRCode
	for rows.Next() {
		var item models.QRCode
		if err := rows.Scan(&item.ID, &item.Code, &item.OfficeID, &item.DeviceType, &item.DeviceModel, &item.CreatedAt, &item.RoomNumber, &item.IsConfigured); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan qrcode"})
		}
		items = append(items, item)
	}
	return c.JSON(items)
}

func (h *QRCodeHandler) Create(c *fiber.Ctx) error {
	var req models.QRCodeCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Count <= 0 || req.Count > 50 {
		req.Count = 1
	}

	var codes []models.QRCode
	for i := 0; i < req.Count; i++ {
		var item models.QRCode
		code := uuid.New().String()[:12]
		err := h.DB.QueryRow(
			"INSERT INTO qr_codes (code) VALUES ($1) RETURNING id, code, created_at",
			code,
		).Scan(&item.ID, &item.Code, &item.CreatedAt)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create qrcode"})
		}
		codes = append(codes, item)
	}
	return c.Status(201).JSON(codes)
}

func (h *QRCodeHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	var req models.QRCodeUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if req.OfficeID != nil {
		var roomNumber string
		h.DB.QueryRow("SELECT room_number FROM offices WHERE id = $1", *req.OfficeID).Scan(&roomNumber)
		if roomNumber == "" {
			return c.Status(400).JSON(models.ErrorResponse{Error: "office not found"})
		}
	}

	_, err = h.DB.Exec(
		`UPDATE qr_codes SET office_id = $1, device_type = $2, device_model = $3 WHERE id = $4`,
		req.OfficeID, req.DeviceType, req.DeviceModel, id,
	)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update qrcode"})
	}
	return c.JSON(fiber.Map{"message": "qrcode updated"})
}

func (h *QRCodeHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	_, err = h.DB.Exec("DELETE FROM qr_codes WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete qrcode"})
	}
	return c.JSON(fiber.Map{"message": "qrcode deleted"})
}

func (h *QRCodeHandler) Image(c *fiber.Ctx) error {
	code := c.Params("code")

	var officeID *int
	var roomNumber, deviceType, deviceModel string
	h.DB.QueryRow(`SELECT q.office_id, COALESCE(o.room_number, ''), COALESCE(o.device_type, ''), COALESCE(o.device_model, '')
		FROM qr_codes q LEFT JOIN offices o ON q.office_id = o.id WHERE q.code = $1`, code,
	).Scan(&officeID, &roomNumber, &deviceType, &deviceModel)

	serverURL := c.Protocol() + "://" + c.Hostname()
	if c.Port() != "443" && c.Port() != "80" {
		serverURL += ":" + c.Port()
	}
	dataURL := serverURL + "/ticket?office=" + code

	qr, err := qrcode.New(dataURL, qrcode.Medium)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to generate qrcode"})
	}

	img := qr.Image(256)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to encode image"})
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return c.JSON(fiber.Map{"image": "data:image/png;base64," + base64Str, "data_url": dataURL})
}

func (h *QRCodeHandler) DeviceModels(c *fiber.Ctx) error {
	rows, err := h.DB.Query(`SELECT DISTINCT device_model FROM (
		SELECT device_model FROM qr_codes WHERE device_model IS NOT NULL AND device_model != ''
		UNION
		SELECT device_model FROM offices WHERE device_model IS NOT NULL AND device_model != ''
	) AS all_models ORDER BY device_model`)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query device models"})
	}
	defer rows.Close()

	var models_ []string
	for rows.Next() {
		var m string
		if rows.Scan(&m) == nil {
			models_ = append(models_, m)
		}
	}
	if models_ == nil {
		models_ = []string{}
	}
	return c.JSON(fiber.Map{"models": models_})
}
