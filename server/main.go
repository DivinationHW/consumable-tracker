package main

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"consumable-tracker/config"
	"consumable-tracker/database"
	"consumable-tracker/handlers"
	"consumable-tracker/middleware"
)

func main() {
	cfg := config.Load()

	md := handlers.NewWebSocketHub()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	adminHash := sha256.Sum256([]byte(cfg.AdminPassword))
	if err := database.InitAdmin(cfg.AdminUsername, fmt.Sprintf("%x", adminHash)); err != nil {
		log.Printf("Warning: failed to init admin: %v", err)
	}

	middleware.SetJWTSecret(cfg.JWTSecret)

	backupDir := "/backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		backupDir = filepath.Join(".", "backups")
		os.MkdirAll(backupDir, 0755)
	}
	handlers.StartBackupScheduler(database.DB, backupDir)

	authH := handlers.NewAuthHandler(database.DB)
	userH := handlers.NewUserHandler(database.DB)
	consH := handlers.NewConsumableHandler(database.DB)
	officeH := handlers.NewOfficeHandler(database.DB)
	recordH := handlers.NewRecordHandler(database.DB)
	statsH := handlers.NewStatsHandler(database.DB)
	noteH := handlers.NewNoteHandler(database.DB)
	ptH := handlers.NewProblemTypeHandler(database.DB)
	qrH := handlers.NewQRCodeHandler(database.DB)
	ticketH := handlers.NewTicketHandler(database.DB, md)
	backupH := handlers.NewBackupHandler(database.DB, backupDir)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.HSTS())

	api := app.Group("/api")

	api.Post("/auth/login", middleware.PublicRateLimit(), authH.Login)
	api.Put("/auth/password", middleware.AuthRequired(), authH.ChangePassword)

	adminAPI := api.Group("", middleware.AuthRequired(), middleware.RequireAdmin())
	adminAPI.Get("/users", userH.List)
	adminAPI.Post("/users", userH.Create)
	adminAPI.Put("/users/:id", userH.Update)
	adminAPI.Delete("/users/:id", userH.Delete)

	roAPI := api.Group("", middleware.AuthRequired(), middleware.AdminOrReadonly())
	roAPI.Get("/consumables", consH.List)
	roAPI.Get("/offices", officeH.List)
	roAPI.Get("/records", recordH.List)
	roAPI.Get("/stats/summary", statsH.Summary)
	roAPI.Get("/problem-types", ptH.List)
	roAPI.Get("/tickets", ticketH.List)

	adminAPI.Post("/consumables", consH.Create)
	adminAPI.Put("/consumables/:id", consH.Update)
	adminAPI.Delete("/consumables/:id", consH.Delete)

	adminAPI.Post("/offices", officeH.Create)
	adminAPI.Put("/offices/:id", officeH.Update)
	adminAPI.Delete("/offices/:id", officeH.Delete)

	adminAPI.Post("/records", recordH.Create)
	adminAPI.Put("/records/:id", recordH.Update)
	adminAPI.Delete("/records/:id", recordH.Delete)

	adminAPI.Get("/notes", noteH.List)
	adminAPI.Post("/notes", noteH.Create)
	adminAPI.Put("/notes/:id", noteH.Update)
	adminAPI.Delete("/notes/:id", noteH.Delete)

	adminAPI.Post("/problem-types", ptH.Create)
	adminAPI.Put("/problem-types/:id", ptH.Update)
	adminAPI.Delete("/problem-types/:id", ptH.Delete)
	adminAPI.Put("/problem-types/sort", ptH.Sort)

	adminAPI.Post("/qrcodes", middleware.QRCodeRateLimit(), qrH.Create)
	adminAPI.Get("/qrcodes", qrH.List)
	adminAPI.Put("/qrcodes/:id", qrH.Update)
	adminAPI.Delete("/qrcodes/:id", qrH.Delete)
	adminAPI.Get("/device-models", qrH.DeviceModels)

	adminAPI.Get("/backup/config", backupH.GetConfig)
	adminAPI.Put("/backup/config", backupH.SaveConfig)
	adminAPI.Get("/backup/list", backupH.List)
	adminAPI.Post("/backup/now", backupH.CreateNow)
	adminAPI.Get("/backup/download/:filename", backupH.Download)
	adminAPI.Post("/backup/restore/:filename", backupH.Restore)
	adminAPI.Delete("/backup/:filename", backupH.Delete)

	adminAPI.Post("/tickets/:id/status", ticketH.UpdateStatus)
	adminAPI.Post("/tickets/:id/complete", ticketH.Complete)
	adminAPI.Delete("/tickets/:id", ticketH.Delete)

	ticketAPI := app.Group("/ticket")
	ticketAPI.Post("/", middleware.PublicRateLimit(), ticketH.Submit)
	ticketAPI.Get("/:id", ticketH.GetStatus)
	ticketAPI.Get("/office/:office_id/problem-types", ticketH.OfficeProblemTypes)

	adminAPI.Get("/qrcodes/:code/image", qrH.Image)

	app.Static("/", "./web-dist")

	handlers.SetupWebSocketRoute(app, md)

	port := ":" + cfg.Port
	log.Printf("Starting server on port %s (HTTPS mode: %s)", port, cfg.HTTPSMode)

	var err error
	if cfg.HTTPSMode == "http" {
		err = app.Listen(port)
	} else {
		certFile := resolveCertPath("fullchain.pem")
		keyFile := resolveCertPath("privkey.pem")

		if _, cerr := os.Stat(certFile); os.IsNotExist(cerr) {
			log.Printf("Certificate not found at %s, generating self-signed...", certFile)
			generateSelfSignedCert(certFile, keyFile, cfg.Domain)
		}

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		err = app.Listen(port, fiber.ListenConfig{
			CertFile:          certFile,
			CertKeyFile:       keyFile,
			TLSConfig:         tlsConfig,
		})
	}

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func resolveCertPath(filename string) string {
	paths := []string{
		"/certs/" + filename,
		"./certs/" + filename,
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "./certs/" + filename
}

func generateSelfSignedCert(certFile, keyFile, domain string) {
	os.MkdirAll(filepath.Dir(certFile), 0755)
	cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-days", "365",
		"-newkey", "rsa:2048",
		"-keyout", keyFile,
		"-out", certFile,
		"-subj", fmt.Sprintf("/CN=%s", domain))
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: failed to generate self-signed cert: %v", err)
	}
	os.Chmod(keyFile, 0600)
	os.Chmod(certFile, 0644)
}
