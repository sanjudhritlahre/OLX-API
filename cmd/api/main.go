package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjudhritlahre/olx-api/internal/config"
	"github.com/sanjudhritlahre/olx-api/internal/db"
	"github.com/sanjudhritlahre/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// Ensure Proper Cleanup
	defer db.Close()

	fmt.Println("Database Connected.")
	fmt.Println("OLX-API server is running...")

	// Database Initialization
	lh := handlers.NewListingHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.Listings)
	mux.HandleFunc("DELETE /listings/{id}", lh.DeleteListing)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// Listen for termination signals and gracefully shut down the server.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigChan)

		sig := <-sigChan
		slog.Info("Shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Shutting down HTTP server", "timeout", "10s")

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("HTTP server shutdown failed", "error", err)
			return
		}

		slog.Info("HTTP server shutdown completed successfully.")
	}()

	slog.Info("HTTP server starting", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("HTTP server stopped")
}
