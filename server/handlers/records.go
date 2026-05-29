package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"consumable-tracker/models"
)

type RecordHandler struct {
	DB *sql.DB
}

func (h *RecordHandler) List(c *fiber.Ctx) error {
	officeID := c.Query("office_id")
	consumableID := c.Query("consumable_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 50)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 50
	}

	where := []string{}
	args := []interface{}{}
	if officeID != "" {
		where = append(where, "r.office_id = ?")
		args = append(args, officeID)
	}
	if consumableID != "" {
		where = append(where, "r.consumable_id = ?")
		args = append(args, consumableID)
	}
	if dateFrom != "" {
		where = append(where, "r.usage_date >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		where = append(where, "r.usage_date <= ?")
		args = append(args, dateTo)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM usage_records r %s", whereClause)
	h.DB.QueryRow(countQ, args...).Scan(&total)

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT r.id, r.office_id, r.consumable_id, r.quantity, r.usage_date, r.note,
		       r.created_at, r.updated_at,
		       o.room_number, c.name, o.device_type, o.device_model
		FROM usage_records r
		JOIN offices o ON r.office_id = o.id
		JOIN consumable_models c ON r.consumable_id = c.id
		%s
		ORDER BY r.usage_date DESC, r.created_at DESC
		LIMIT ? OFFSET ?`, whereClause)
	args = append(args, pageSize, offset)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败", "detail": err.Error()})
	}
	defer rows.Close()

	list := []models.UsageRecord{}
	for rows.Next() {
		var r models.UsageRecord
		if err := rows.Scan(&r.ID, &r.OfficeID, &r.ConsumableID, &r.Quantity, &r.UsageDate, &r.Note,
			&r.CreatedAt, &r.UpdatedAt, &r.OfficeName, &r.ConsumableName, &r.DeviceType, &r.DeviceModel); err != nil {
			continue
		}
		list = append(list, r)
	}
	return c.JSON(fiber.Map{
		"data":  list,
		"total": total,
		"page":  page,
	})
}

func (h *RecordHandler) Create(c *fiber.Ctx) error {
	var r models.UsageRecord
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if r.OfficeID == 0 || r.ConsumableID == 0 || r.Quantity < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "参数不完�?})
	}
	if r.UsageDate == "" {
		r.UsageDate = strings.SplitN(c.Get("Date", "now"), "T", 2)[0]
	}
	result, err := h.DB.Exec("INSERT INTO usage_records (office_id, consumable_id, quantity, usage_date, note) VALUES (?, ?, ?, ?, ?)",
		r.OfficeID, r.ConsumableID, r.Quantity, r.UsageDate, r.Note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建失败"})
	}
	id, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *RecordHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	var r models.UsageRecord
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if r.Quantity < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "数量必须大于0"})
	}
	_, err = h.DB.Exec("UPDATE usage_records SET office_id=?, consumable_id=?, quantity=?, usage_date=?, note=?, updated_at=datetime('now') WHERE id=?",
		r.OfficeID, r.ConsumableID, r.Quantity, r.UsageDate, r.Note, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}
	return c.JSON(fiber.Map{"message": "更新成功"})
}

func (h *RecordHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效ID"})
	}
	_, err = h.DB.Exec("DELETE FROM usage_records WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

func (h *RecordHandler) Export(c *fiber.Ctx) error {
	officeID := c.Query("office_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	where := []string{}
	args := []interface{}{}
	if officeID != "" {
		where = append(where, "r.office_id = ?")
		args = append(args, officeID)
	}
	if dateFrom != "" {
		where = append(where, "r.usage_date >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		where = append(where, "r.usage_date <= ?")
		args = append(args, dateTo)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT r.usage_date, o.room_number, o.device_type, o.device_model,
		       c.name, c.unit, r.quantity, r.note
		FROM usage_records r
		JOIN offices o ON r.office_id = o.id
		JOIN consumable_models c ON r.consumable_id = c.id
		%s
		ORDER BY r.usage_date DESC, r.created_at DESC`, whereClause)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "导出失败"})
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "日期")
	f.SetCellValue(sheet, "B1", "房间�?)
	f.SetCellValue(sheet, "C1", "设备类型")
	f.SetCellValue(sheet, "D1", "设备型号")
	f.SetCellValue(sheet, "E1", "耗材名称")
	f.SetCellValue(sheet, "F1", "单位")
	f.SetCellValue(sheet, "G1", "数量")
	f.SetCellValue(sheet, "H1", "备注")

	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheet, "A1", "H1", style)

	idx := 2
	for rows.Next() {
		var date, room, dtype, dmodel, cname, unit string
		var qty int
		var note string
		if err := rows.Scan(&date, &room, &dtype, &dmodel, &cname, &unit, &qty, &note); err != nil {
			continue
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", idx), date)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", idx), room)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", idx), dtype)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", idx), dmodel)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", idx), cname)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", idx), unit)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", idx), qty)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", idx), note)
		idx++
	}

	f.SetColWidth(sheet, "A", 12)
	f.SetColWidth(sheet, "B", 12)
	f.SetColWidth(sheet, "C", 12)
	f.SetColWidth(sheet, "D", 20)
	f.SetColWidth(sheet, "E", 25)
	f.SetColWidth(sheet, "F", 8)
	f.SetColWidth(sheet, "G", 8)
	f.SetColWidth(sheet, "H", 20)

	c.Response().Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header("Content-Disposition", "attachment; filename=records.xlsx")
	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "写入失败"})
	}
	return nil
}
