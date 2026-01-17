package main

import (
	"log"
	"net/http"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/config"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/handler"


)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
    mux.HandleFunc("/health", handler.Health)

	log.Printf("%s running on :%s", cfg.ServiceName, cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}