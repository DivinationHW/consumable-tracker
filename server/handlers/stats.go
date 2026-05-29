package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/server/models"
)

type StatsHandler struct {
	DB *sql.DB
}

func (h *StatsHandler) Summary(c *fiber.Ctx) error {
	var total int
	h.DB.QueryRow("SELECT COALESCE(SUM(quantity),0) FROM usage_records").Scan(&total)

	var month int
	h.DB.QueryRow(`SELECT COALESCE(SUM(quantity),0) FROM usage_records
		WHERE strftime('%Y-%m', usage_date) = strftime('%Y-%m', 'now')`).Scan(&month)

	var topOffice string
	h.DB.QueryRow(`SELECT o.room_number FROM usage_records r
		JOIN offices o ON r.office_id = o.id
		GROUP BY o.room_number ORDER BY SUM(r.quantity) DESC LIMIT 1`).Scan(&topOffice)

	var topConsumable string
	h.DB.QueryRow(`SELECT c.name FROM usage_records r
		JOIN consumable_models c ON r.consumable_id = c.id
		GROUP BY c.name ORDER BY SUM(r.quantity) DESC LIMIT 1`).Scan(&topConsumable)

	byOffice := []models.StatsItem{}
	rows, err := h.DB.Query(`SELECT o.room_number, COALESCE(SUM(r.quantity),0)
		FROM usage_records r JOIN offices o ON r.office_id = o.id
		GROUP BY o.room_number ORDER BY SUM(r.quantity) DESC`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item models.StatsItem
			if err := rows.Scan(&item.Label, &item.Value); err == nil {
				byOffice = append(byOffice, item)
			}
		}
	}

	trend := []models.StatsItem{}
	trows, err := h.DB.Query(`SELECT strftime('%Y-%m', usage_date) AS m, COALESCE(SUM(quantity),0)
		FROM usage_records GROUP BY m ORDER BY m DESC LIMIT 12`)
	if err == nil {
		defer trows.Close()
		for trows.Next() {
			var item models.StatsItem
			if err := trows.Scan(&item.Label, &item.Value); err == nil {
				trend = append(trend, item)
			}
		}
	}

	return c.JSON(models.StatsSummary{
		TotalUsage:     total,
		CurrentMonth:   month,
		TopOffice:      topOffice,
		TopConsumable:  topConsumable,
		ByOffice:       byOffice,
		MonthlyTrend:   trend,
	})
}
