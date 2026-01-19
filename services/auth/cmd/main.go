package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/handler"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/middleware"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/service"
	
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("failed to load:", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)

	// Initialize services and handlers
	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	// public
	r.Get("/health", handler.Health)
	r.Post("/login", authHandler.Login)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		// endpoint protected nanti di sini
	})

	log.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}