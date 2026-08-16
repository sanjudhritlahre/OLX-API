package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sanjudhritlahre/olx-api/internal/config"
	"github.com/sanjudhritlahre/olx-api/internal/db"
	"github.com/sanjudhritlahre/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	
	_, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	fmt.Println("Database Connected.")
	fmt.Println("OLX-API server is running...")
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Failed: %v", err)
	}
}