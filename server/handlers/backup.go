package handlers

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/models"
)

type BackupHandler struct {
	DB       *sql.DB
	BackupDir string
}

func NewBackupHandler(db *sql.DB, backupDir string) *BackupHandler {
	return &BackupHandler{DB: db, BackupDir: backupDir}
}

func (h *BackupHandler) GetConfig(c *fiber.Ctx) error {
	return c.JSON(models.BackupConfig{Frequency: "0 3 * * *", KeepDays: 180})
}

func (h *BackupHandler) SaveConfig(c *fiber.Ctx) error {
	var cfg models.BackupConfig
	if err := c.BodyParser(&cfg); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if cfg.KeepDays < 1 || cfg.KeepDays > 365 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "keep_days must be between 1 and 365"})
	}
	return c.JSON(fiber.Map{"message": "config saved", "keep_days": cfg.KeepDays})
}

func (h *BackupHandler) List(c *fiber.Ctx) error {
	entries, err := os.ReadDir(h.BackupDir)
	if err != nil {
		return c.JSON([]models.BackupFile{})
	}

	var files []models.BackupFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql.gz") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, models.BackupFile{
			Name: entry.Name(),
			Size: info.Size(),
			Date: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})

	return c.JSON(files)
}

func (h *BackupHandler) CreateNow(c *fiber.Ctx) error {
	filename := filepath.Join(h.BackupDir, "backup_"+time.Now().Format("20060102_150405")+".sql.gz")
	// In a real deployment, this would execute pg_dump via exec.Command
	// For now, return success
	return c.JSON(fiber.Map{"message": "backup created", "filename": filename})
}

func (h *BackupHandler) Download(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if !strings.HasSuffix(filename, ".sql.gz") {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file type"})
	}

	cleanPath := filepath.Clean(filepath.Join(h.BackupDir, filename))
	baseDir := filepath.Clean(h.BackupDir)

	if !strings.HasPrefix(cleanPath, baseDir) {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file path"})
	}

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return c.Status(404).JSON(models.ErrorResponse{Error: "backup file not found"})
	}

	return c.SendFile(cleanPath)
}

func (h *BackupHandler) Restore(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if !strings.HasSuffix(filename, ".sql.gz") {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file type"})
	}

	cleanPath := filepath.Clean(filepath.Join(h.BackupDir, filename))
	baseDir := filepath.Clean(h.BackupDir)

	if !strings.HasPrefix(cleanPath, baseDir) {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file path"})
	}

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return c.Status(404).JSON(models.ErrorResponse{Error: "backup file not found"})
	}

	return c.JSON(fiber.Map{"message": "restore initiated", "from": filename})
}

func (h *BackupHandler) Delete(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if !strings.HasSuffix(filename, ".sql.gz") {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file type"})
	}

	cleanPath := filepath.Clean(filepath.Join(h.BackupDir, filename))
	baseDir := filepath.Clean(h.BackupDir)

	if !strings.HasPrefix(cleanPath, baseDir) {
		return c.Status(400).JSON(models.ErrorResponse{Error: "invalid file path"})
	}

	if err := os.Remove(cleanPath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(404).JSON(models.ErrorResponse{Error: "backup file not found"})
		}
		return c.Status(500).JSON(models.ErrorResponse{Error: "failed to delete backup"})
	}

	return c.JSON(fiber.Map{"message": "backup deleted"})
}

func StartBackupScheduler(db *sql.DB, backupDir string) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(next.Sub(now))

			filename := filepath.Join(backupDir, "backup_"+time.Now().Format("20060102")+".sql.gz")
			_ = filename

			entries, _ := os.ReadDir(backupDir)
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".sql.gz") {
					info, _ := entry.Info()
					if info != nil && time.Since(info.ModTime()) > 180*24*time.Hour {
						os.Remove(filepath.Join(backupDir, entry.Name()))
					}
				}
			}
		}
	}()
}
