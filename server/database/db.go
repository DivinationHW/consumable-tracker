package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	for i := 0; i < 30; i++ {
		err = DB.Ping()
		if err == nil {
			log.Println("Database connected successfully")
			return nil
		}
		log.Printf("Waiting for database... attempt %d/30: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to database after 30 attempts: %w", err)
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func InitAdmin(username, passwordHash string) error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = DB.Exec(
			"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, 'admin')",
			username, passwordHash,
		)
		if err != nil {
			return err
		}
		log.Printf("Admin user '%s' created", username)
	}
	return nil
}
