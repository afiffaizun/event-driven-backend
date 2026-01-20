package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/db"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/handler"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/middleware"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/service"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// root context
	ctx := context.Background()

	// database
	dbPool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// repository
	userRepo := repository.NewPostgresUserRepository(dbPool)

	// service
	authService := service.NewAuthService(userRepo)

	// handler
	authHandler := handler.NewAuthHandler(authService)

	// router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.Health)
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			// protected endpoints later
		})
	})

	// http server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// run server
	go func() {
		log.Printf("🚀 %s running on port %s", cfg.AppName, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🛑 shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("✅ server exited properly")
}
