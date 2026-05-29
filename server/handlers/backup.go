package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"consumable-tracker/server/models"
	"modernc.org/sqlite"
)

type BackupHandler struct {
	DB      *sql.DB
	DataDir string
}

func (h *BackupHandler) List(c *fiber.Ctx) error {
	backupsDir := filepath.Join(h.DataDir, "../backups")
	absDir, err := filepath.Abs(backupsDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "服务器错�?})
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		if os.IsNotExist(err) {
			return c.JSON([]models.BackupInfo{})
		}
		return c.Status(500).JSON(fiber.Map{"error": "读取备份目录失败"})
	}

	list := []models.BackupInfo{}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".db" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		list = append(list, models.BackupInfo{
			Filename:  e.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}
	return c.JSON(list)
}

func (h *BackupHandler) Create(c *fiber.Ctx) error {
	backupsDir := filepath.Join(h.DataDir, "../backups")
	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建备份目录失败"})
	}

	filename := fmt.Sprintf("backup_%s.db", time.Now().Format("20060102_150405"))
	dstPath := filepath.Join(backupsDir, filename)
	db := h.DB

	conn, err := db.Conn(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取数据库连接失�?})
	}
	defer conn.Close()

	err = conn.Raw(func(driverConn interface{}) error {
		sqliteConn, ok := driverConn.(*sqlite.Conn)
		if !ok {
			return fmt.Errorf("not a sqlite conn")
		}
		dstDB, err := sql.Open("sqlite", dstPath)
		if err != nil {
			return err
		}
		defer dstDB.Close()

		dstConn, err := dstDB.Conn(c.Context())
		if err != nil {
			return err
		}
		defer dstConn.Close()

		return dstConn.Raw(func(dstDriver interface{}) error {
			dstSQLite, ok := dstDriver.(*sqlite.Conn)
			if !ok {
				return fmt.Errorf("dst not sqlite")
			}
			_, err := sqliteConn.Backup("main", dstSQLite, "main")
			return err
		})
	})

	if err != nil {
		log.Printf("[Backup] failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "备份失败"})
	}

	log.Printf("[Backup] created: %s", filename)
	return c.JSON(fiber.Map{"message": "备份成功", "filename": filename})
}

func (h *BackupHandler) Restore(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" || filepath.Ext(filename) != ".db" {
		return c.Status(400).JSON(fiber.Map{"error": "无效文件�?})
	}
	backupsDir := filepath.Join(h.DataDir, "../backups")
	srcPath := filepath.Join(backupsDir, filename)

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{"error": "备份文件不存�?})
	}

	dbPath := filepath.Join(h.DataDir, "data.db")
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "读取备份失败"})
	}
	if err := os.WriteFile(dbPath, srcData, 0644); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "恢复失败"})
	}

	log.Printf("[Backup] restored from: %s (restart required)", filename)
	return c.JSON(fiber.Map{"message": "恢复成功，请重启服务�?})
}

func (h *BackupHandler) Delete(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" || filepath.Ext(filename) != ".db" {
		return c.Status(400).JSON(fiber.Map{"error": "无效文件�?})
	}
	backupsDir := filepath.Join(h.DataDir, "../backups")
	path := filepath.Join(backupsDir, filename)
	if err := os.Remove(path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除失败"})
	}
	return c.JSON(fiber.Map{"message": "删除成功"})
}

func (h *BackupHandler) StartAutoBackup() {
	backupsDir := filepath.Join(h.DataDir, "../backups")
	os.MkdirAll(backupsDir, 0755)

	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(next.Sub(now))

			filename := fmt.Sprintf("auto_%s.db", time.Now().Format("20060102"))
			dstPath := filepath.Join(backupsDir, filename)
			if _, err := os.Stat(dstPath); err == nil {
				continue
			}

			db := h.DB
			conn, err := db.Conn(nil)
			if err != nil {
				log.Printf("[AutoBackup] conn error: %v", err)
				continue
			}

			err = conn.Raw(func(driverConn interface{}) error {
				sqliteConn, ok := driverConn.(*sqlite.Conn)
				if !ok {
					return fmt.Errorf("not a sqlite conn")
				}
				dstDB, err := sql.Open("sqlite", dstPath)
				if err != nil {
					return err
				}
				defer dstDB.Close()

				dstConn, err := dstDB.Conn(nil)
				if err != nil {
					return err
				}
				defer dstConn.Close()

				return dstConn.Raw(func(dstDriver interface{}) error {
					dstSQLite, ok := dstDriver.(*sqlite.Conn)
					if !ok {
						return fmt.Errorf("dst not sqlite")
					}
					_, err := sqliteConn.Backup("main", dstSQLite, "main")
					return err
				})
			})
			conn.Close()

			if err != nil {
				log.Printf("[AutoBackup] failed: %v", err)
				continue
			}
			log.Printf("[AutoBackup] created: %s", filename)

			entries, _ := filepath.Glob(filepath.Join(backupsDir, "auto_*.db"))
			if len(entries) > 180 {
				for _, e := range entries[:len(entries)-180] {
					os.Remove(e)
				}
			}
		}
	}()
}
