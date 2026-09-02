package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/sanjudhritlahre/olx-api/internal/httpx"
	"github.com/sanjudhritlahre/olx-api/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db: db,
		logger: logger,
	}
}

func (lh ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`)
	if err != nil {
		lh.logger.Info("lh.db.QueryContext error", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Unable to retrieve listings at this time.", httpx.InternalError)
		return
	}
	defer rows.Close()

	listings := []listing{}
	for rows.Next() {
		var l listing

		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Info("rows.Scan error", "err", err)
			httpx.Error(w, http.StatusInternalServerError, "Unable to process listing data.", httpx.InternalError)
			return
		}

		lh.logger.Info("listing fetched", "total", len(listings))
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.Err: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Unable to complete listing retrieval.", httpx.InternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete failed!", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Unable to delete the listing at this time.", httpx.InternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}