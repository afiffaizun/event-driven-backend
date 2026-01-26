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
	cfg, _ := config.Load()
	ctx := context.Background()

	dbPool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	userRepo := repository.NewPostgresUserRepository(dbPool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))
			r.Get("/protected", func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("protected content"))
			})
		})
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go server.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctxShutdown)
}
