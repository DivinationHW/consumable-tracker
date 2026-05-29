package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type StatsHandler struct {
	DB *sql.DB
}

func NewStatsHandler(db *sql.DB) *StatsHandler {
	return &StatsHandler{DB: db}
}

func (h *StatsHandler) Summary(c *fiber.Ctx) error {
	var stats models.StatsSummary

	h.DB.QueryRow("SELECT COALESCE(SUM(quantity), 0) FROM usage_records").Scan(&stats.TotalUsage)
	h.DB.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM usage_records 
		WHERE to_char(usage_date, 'YYYY-MM') = to_char(NOW(), 'YYYY-MM')`).Scan(&stats.CurrentMonthUsage)
	h.DB.QueryRow(`SELECT COALESCE(MIN(usage_date)::text, ''), COALESCE(MAX(usage_date)::text, '') 
		FROM usage_records`).Scan(&stats.DateRangeStart, &stats.DateRangeEnd)

	h.DB.QueryRow(`SELECT o.room_number, SUM(r.quantity) as total
		FROM usage_records r JOIN offices o ON r.office_id = o.id
		GROUP BY o.room_number ORDER BY total DESC LIMIT 1`).Scan(&stats.MostUsedOffice, &stats.MostUsedOfficeCount)

	officeRows, err := h.DB.Query(`SELECT o.id, o.room_number, COALESCE(SUM(r.quantity), 0) as total
		FROM offices o LEFT JOIN usage_records r ON o.id = r.office_id
		GROUP BY o.id, o.room_number ORDER BY total DESC`)
	if err == nil {
		defer officeRows.Close()
		for officeRows.Next() {
			var os models.OfficeStat
			if officeRows.Scan(&os.OfficeID, &os.RoomNumber, &os.Total) == nil {
				stats.OfficeStats = append(stats.OfficeStats, os)
			}
		}
	}

	monthlyRows, err := h.DB.Query(`SELECT to_char(usage_date, 'YYYY-MM') as month, 
		EXTRACT(YEAR FROM usage_date)::int as year, SUM(quantity) as total
		FROM usage_records GROUP BY month, year ORDER BY month DESC LIMIT 12`)
	if err == nil {
		defer monthlyRows.Close()
		for monthlyRows.Next() {
			var ms models.MonthlyStat
			if monthlyRows.Scan(&ms.Month, &ms.Year, &ms.Total) == nil {
				stats.MonthlyStats = append(stats.MonthlyStats, ms)
			}
		}
	}

	if stats.OfficeStats == nil {
		stats.OfficeStats = []models.OfficeStat{}
	}
	if stats.MonthlyStats == nil {
		stats.MonthlyStats = []models.MonthlyStat{}
	}

	return c.JSON(stats)
}
