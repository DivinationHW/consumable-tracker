package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"golang.org/x/crypto/bcrypt"

	"consumable-tracker/database"
	"consumable-tracker/handlers"
	"consumable-tracker/middleware"
)

//go:embed web-dist/*
var webDist embed.FS

func main() {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/app/data"
	}

	db, err := database.NewDB(dataDir)
	if err != nil {
		log.Fatalf("[FATAL] database: %v", err)
	}
	defer db.Close()

	initDefaultAdmin(db)

	wsHub := handlers.NewWSHub()

	authHandler := &handlers.AuthHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db}
	consumableHandler := &handlers.ConsumableHandler{DB: db}
	officeHandler := &handlers.OfficeHandler{DB: db}
	recordHandler := &handlers.RecordHandler{DB: db}
	statsHandler := &handlers.StatsHandler{DB: db}
	noteHandler := &handlers.NoteHandler{DB: db}
	problemTypeHandler := &handlers.ProblemTypeHandler{DB: db}
	qrHandler := &handlers.QRCodeHandler{DB: db}
	ticketHandler := &handlers.TicketHandler{DB: db, Broadcast: handlers.BroadcastFunc(wsHub)}
	backupHandler := &handlers.BackupHandler{DB: db, DataDir: dataDir}
	backupHandler.StartAutoBackup()

	app := fiber.New(fiber.Config{
		AppName:     "ConsumableTracker",
		BodyLimit:   10 * 1024 * 1024,
		ReadTimeout: 30,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())

	api := app.Group("/api")

	api.Post("/login", middleware.LoginRateLimit(), authHandler.Login)

	api.Use(middleware.JWTAuth())

	api.Get("/me", authHandler.Me)
	api.Post("/change-password", authHandler.ChangePassword)

	admin := api.Group("", middleware.RoleAdmin())

	admin.Get("/users", userHandler.List)
	admin.Post("/users", userHandler.Create)
	admin.Put("/users/:id", userHandler.Update)
	admin.Delete("/users/:id", userHandler.Delete)

	admin.Get("/consumables", consumableHandler.List)
	admin.Post("/consumables", consumableHandler.Create)
	admin.Put("/consumables/:id", consumableHandler.Update)
	admin.Delete("/consumables/:id", consumableHandler.Delete)

	admin.Get("/offices", officeHandler.List)
	admin.Post("/offices", officeHandler.Create)
	admin.Put("/offices/:id", officeHandler.Update)
	admin.Delete("/offices/:id", officeHandler.Delete)

	admin.Get("/notes", noteHandler.List)
	admin.Post("/notes", noteHandler.Create)
	admin.Put("/notes/:id", noteHandler.Update)
	admin.Delete("/notes/:id", noteHandler.Delete)

	admin.Get("/problem-types", problemTypeHandler.List)
	admin.Post("/problem-types", problemTypeHandler.Create)
	admin.Put("/problem-types/:id", problemTypeHandler.Update)
	admin.Delete("/problem-types/:id", problemTypeHandler.Delete)

	admin.Get("/qrcodes", qrHandler.List)
	admin.Post("/qrcodes", qrHandler.Create)
	admin.Post("/qrcodes/bulk", qrHandler.GenerateBulk)
	admin.Delete("/qrcodes/:id", qrHandler.Delete)
	admin.Get("/qrcodes/print", qrHandler.PrintPage)

	admin.Post("/records", recordHandler.Create)
	admin.Put("/records/:id", recordHandler.Update)
	admin.Delete("/records/:id", recordHandler.Delete)

	admin.Get("/backups", backupHandler.List)
	admin.Post("/backups", backupHandler.Create)
	admin.Post("/backups/:filename/restore", backupHandler.Restore)
	admin.Delete("/backups/:filename", backupHandler.Delete)

	readonly := api.Group("")
	readonly.Get("/records", recordHandler.List)
	readonly.Get("/records/export", recordHandler.Export)
	readonly.Get("/stats", statsHandler.Summary)
	readonly.Get("/device-models", officeHandler.GetDeviceModels)
	readonly.Get("/device-types", problemTypeHandler.GetDeviceTypes)

	readonly.Get("/tickets", ticketHandler.List)
	readonly.Get("/tickets/:id", ticketHandler.GetByID)
	readonly.Post("/tickets/:id/process", ticketHandler.Process)

	public := app.Group("/public")
	public.Get("/ticket/:id", middleware.PublicRateLimit(), ticketHandler.GetPublicByID)
	public.Post("/ticket", middleware.PublicRateLimit(), ticketHandler.CreatePublic)

	apiPublic := api.Group("/public")
	apiPublic.Get("/ticket/:id", ticketHandler.GetPublicByID)
	apiPublic.Post("/ticket", ticketHandler.CreatePublic)

	readonly.Get("/qrcodes/:code/image", qrHandler.Image)
	readonly.Get("/qrcodes/:code/image-base64", qrHandler.ImageBase64)

	app.Get("/ws", middleware.OptionalAuth(), handlers.WSUpgrade, websocket.New(wsHub.HandleWS))

	subFS, err := fs.Sub(webDist, "web-dist")
	if err != nil {
		log.Fatalf("[FATAL] static files: %v", err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         subFS,
		Index:        "index.html",
		NotFoundFile: "index.html",
		Browse:       false,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8443"
	}
	addr := ":" + port

	httpsMode := os.Getenv("HTTPS_MODE")
	switch httpsMode {
	case "selfsigned":
		certDir := os.Getenv("CERT_DIR")
		if certDir == "" {
			certDir = "/app/certs"
		}
		certFile := filepath.Join(certDir, "cert.pem")
		keyFile := filepath.Join(certDir, "key.pem")
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Println("[TLS] Generating self-signed certificate...")
			generateSelfSigned(certFile, keyFile)
		}
		log.Printf("[Server] HTTPS on %s (self-signed)", addr)
		log.Fatal(app.ListenTLS(addr, certFile, keyFile))
	case "letsencrypt":
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			log.Fatal("[FATAL] DOMAIN required for letsencrypt mode")
		}
		log.Printf("[Server] HTTPS on %s (Let's Encrypt: %s)", addr, domain)
		log.Fatal(app.Listen(addr))
	default:
		log.Printf("[Server] HTTP on %s", addr)
		log.Fatal(app.Listen(addr))
	}
}

func initDefaultAdmin(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 12)
		db.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, 'admin')", "admin", string(hash))
		log.Println("[Init] Created default admin/admin123")
	}
}

func generateSelfSigned(certFile, keyFile string) {
	os.MkdirAll(filepath.Dir(certFile), 0755)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("[TLS] generate key: %v", err)
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "Consumable Tracker"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:     []string{"localhost"},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("[TLS] create cert: %v", err)
	}
	if err := os.WriteFile(certFile, certDER, 0644); err != nil {
		log.Fatalf("[TLS] write cert: %v", err)
	}
	privBytes, _ := x509.MarshalECPrivateKey(priv)
	if err := os.WriteFile(keyFile, privBytes, 0600); err != nil {
		log.Fatalf("[TLS] write key: %v", err)
	}
}
