package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/db"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

func main() {
	// Use localhost for running outside Docker
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://authuser:authpass@localhost:5432/authdb?sslmode=disable")
	}
	
	cfg, _ := config.Load()
	ctx := context.Background()

	dbPool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// Hash password "admin"
	hashedPassword, err := security.HashPassword("admin")
	if err != nil {
		log.Fatal(err)
	}

	// Create admin user
	_, err = dbPool.Exec(ctx, `
		INSERT INTO users (username, password) 
		VALUES ($1, $2)
		ON CONFLICT (username) DO UPDATE SET password = $2
	`, "admin", hashedPassword)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin user created successfully!")
	fmt.Println("Username: admin")
	fmt.Println("Password: admin")
}
