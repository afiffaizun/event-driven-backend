package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/db"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/handler"
	authMiddleware "github.com/afiffaizun/event-driven-backend/services/auth/internal/middleware"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	dbPool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	userRepo := repository.NewPostgresUserRepository(dbPool)
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(dbPool)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)

	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler(dbPool)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(authMiddleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Auth(cfg.JWTSecret))
			r.Get("/protected", func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("protected content"))
			})
		})
	})

	r.Get("/health", healthHandler.Check)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting %s on port %s", cfg.AppName, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
