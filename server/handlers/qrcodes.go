package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"consumable-tracker/models"
)

type QRCodeHandler struct {
	DB *sql.DB
}

func (h *QRCodeHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Query(`SELECT q.id, q.code, q.office_id, q.device_type, q.device_model, q.created_at,
		COALESCE(o.room_number,'') FROM qr_codes q LEFT JOIN offices o ON q.office_id = o.id ORDER BY q.created_at DESC`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()
	list := []models.QRCode{}
	for rows.Next() {
		var q models.QRCode
		if err := rows.Scan(&q.ID, &q.Code, &q.OfficeID, &q.DeviceType, &q.DeviceModel, &q.CreatedAt, &q.OfficeName); err != nil {
			continue
		}
		list = append(list, q)
	}
	return c.JSON(list)
}

type createQRReq struct {
	Code        string `json:"code"`
	OfficeID    *int   `json:"office_id"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
}

func (h *QRCodeHandler) Create(c *fiber.Ctx) error {
	var req createQRReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Code == "" {
		req.Code = uuid.New().String()[:8]
	}
	if req.OfficeID != nil && *req.OfficeID == 0 {
		req.OfficeID = nil
	}
	_, err := h.DB.Exec("INSERT INTO qr_codes (code, office_id, device_type, device_model) VALUES (?, ?, ?, ?)",
		req.Code, req.OfficeID, req.DeviceType, req.DeviceModel)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "二维码码值已存在"})
	}
	return c.Status(201).JSON(fiber.Map{"code": req.Code})
}

func (h *QRCodeHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	_, err = h.DB.Exec("DELETE FROM qr_codes WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

func (h *QRCodeHandler) Image(c *fiber.Ctx) error {
	code := c.Params("code")
	baseURL := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	if port := c.Get("X-Forwarded-Port", ""); port != "" {
		baseURL = fmt.Sprintf("%s://%s:%s", c.Protocol(), c.Hostname(), port)
	}
	submitURL := baseURL + "/public/ticket/new?qr=" + url.QueryEscape(code)
	statusURL := baseURL + "/public/ticket/status?qr=" + url.QueryEscape(code)
	content := fmt.Sprintf("耗材使用明细\n报修:%s\n查询:%s", submitURL, statusURL)

	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "生成二维码失�?})
	}
	c.Response().Header("Content-Type", "image/png")
	return c.Send(png)
}

type qrImageResp struct {
	Image string `json:"image"`
}

func (h *QRCodeHandler) ImageBase64(c *fiber.Ctx) error {
	code := c.Params("code")
	baseURL := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	if port := c.Get("X-Forwarded-Port", ""); port != "" {
		baseURL = fmt.Sprintf("%s://%s:%s", c.Protocol(), c.Hostname(), port)
	}
	submitURL := baseURL + "/public/ticket/new?qr=" + url.QueryEscape(code)
	statusURL := baseURL + "/public/ticket/status?qr=" + url.QueryEscape(code)
	content := fmt.Sprintf("耗材使用明细\n报修:%s\n查询:%s", submitURL, statusURL)

	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "生成二维码失�?})
	}
	b64 := base64.StdEncoding.EncodeToString(png)
	return c.JSON(qrImageResp{Image: b64})
}

func (h *QRCodeHandler) GenerateBulk(c *fiber.Ctx) error {
	var req struct {
		Count       int    `json:"count"`
		OfficeID    *int   `json:"office_id"`
		DeviceType  string `json:"device_type"`
		DeviceModel string `json:"device_model"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if req.Count < 1 || req.Count > 100 {
		req.Count = 10
	}
	codes := []string{}
	for i := 0; i < req.Count; i++ {
		code := uuid.New().String()[:8]
		_, err := h.DB.Exec("INSERT INTO qr_codes (code, office_id, device_type, device_model) VALUES (?, ?, ?, ?)",
			code, req.OfficeID, req.DeviceType, req.DeviceModel)
		if err != nil {
			continue
		}
		codes = append(codes, code)
	}
	return c.Status(201).JSON(fiber.Map{"codes": codes, "count": len(codes)})
}

func (h *QRCodeHandler) PrintPage(c *fiber.Ctx) error {
	baseURL := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	if port := c.Get("X-Forwarded-Port", ""); port != "" {
		baseURL = fmt.Sprintf("%s://%s:%s", c.Protocol(), c.Hostname(), port)
	}

	rows, err := h.DB.Query(`SELECT q.code, COALESCE(o.room_number,''), q.device_type, q.device_model
		FROM qr_codes q LEFT JOIN offices o ON q.office_id = o.id ORDER BY q.created_at DESC`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	defer rows.Close()

	var htmlBuf bytes.Buffer
	htmlBuf.WriteString(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>二维码打�?/title>
		<style>body{font-family:sans-serif}.grid{display:grid;grid-template-columns:1fr 1fr 1fr;gap:10px}
		.item{text-align:center;border:1px solid #ddd;padding:8px;page-break-inside:avoid}
		img{width:120px;height:120px}.room{font-weight:bold;margin:4px 0}
		@media print{body{margin:0}.item{border-color:#ccc}}
</style></head><body><div class="grid">`)
	for rows.Next() {
		var code, room, dtype, dmodel string
		if err := rows.Scan(&code, &room, &dtype, &dmodel); err != nil {
			continue
		}
		imgURL := baseURL + "/api/qrcodes/" + code + "/image"
		htmlBuf.WriteString(fmt.Sprintf(`<div class="item">
			<img src="%s" alt="QR" />
			<div class="room">%s</div>
			<div>%s</div>
			<div>%s</div>
		</div>`, imgURL, room, dtype, dmodel))
	}
	htmlBuf.WriteString(`</div></body></html>`)
	c.Response().Header("Content-Type", "text/html; charset=utf-8")
	return c.Send(htmlBuf.Bytes())
}
