package main

import (
	"log"
	"net/http"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/handler"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	r := chi.NewRouter()

	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	r.Get("/health", handler.Health)
	r.Post("/login", authHandler.Login)

	log.Printf("%s running on :%s", cfg.ServiceName, cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
