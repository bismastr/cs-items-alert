package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/services/item"
	"github.com/bismastr/cs-price-alert/internal/services/price"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	db         *db.Db
}

func NewServer(config *config.Config, db *db.Db) (*Server, error) {
	timescaleRepo := timescale_repository.New(db.TimescalePool)
	postgresRepo := repository.New(db.PostgresPool)
	priceService := price.NewPriceService(timescaleRepo, postgresRepo)
	priceHandler := price.NewPriceHandler(priceService)
	itemService := item.NewItemService(timescaleRepo, postgresRepo)
	itemHandler := item.NewItemHandler(itemService)
	router := NewRouter(priceHandler, itemHandler)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Server.Port),
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second, // Keep connection alive for 60s
		MaxHeaderBytes:    1 << 20,          // 1MB max header
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		db:         db,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("Started on port %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return err
}

func (s *Server) Close() error {
	context := context.Background()

	err := s.httpServer.Shutdown(context)
	if err != nil {
		return fmt.Errorf("server failed to close: %w", err)
	}

	s.db.PostgresPool.Close()
	s.db.TimescalePool.Close()

	return nil
}
