package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type NoteHandler struct {
	DB *sql.DB
}

func NewNoteHandler(db *sql.DB) *NoteHandler {
	return &NoteHandler{DB: db}
}

func (h *NoteHandler) List(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	var rows *sql.Rows
	var err error

	if keyword != "" {
		rows, err = h.DB.Query(
			"SELECT id, title, content, created_at, updated_at FROM notes WHERE title ILIKE $1 OR content ILIKE $1 ORDER BY updated_at DESC",
			"%"+keyword+"%",
		)
	} else {
		rows, err = h.DB.Query("SELECT id, title, content, created_at, updated_at FROM notes ORDER BY updated_at DESC")
	}
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to query notes"})
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "failed to scan note"})
		}
		notes = append(notes, n)
	}
	return c.JSON(notes)
}

func (h *NoteHandler) Create(c *fiber.Ctx) error {
	var req models.NoteCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if req.Title == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "title is required"})
	}
	if len(req.Title) > 100 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "title too long (max 100)"})
	}
	if len(req.Content) > 2000 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "content too long (max 2000)"})
	}

	var note models.Note
	err := h.DB.QueryRow(
		"INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at",
		req.Title, req.Content,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to create note"})
	}
	return c.Status(201).JSON(note)
}

func (h *NoteHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.NoteCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	_, err := h.DB.Exec(
		"UPDATE notes SET title = $1, content = $2, updated_at = NOW() WHERE id = $3",
		req.Title, req.Content, id,
	)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to update note"})
	}
	return c.JSON(fiber.Map{"message": "note updated"})
}

func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete note"})
	}
	return c.JSON(fiber.Map{"message": "note deleted"})
}
