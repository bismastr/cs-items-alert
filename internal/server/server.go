package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/services/price"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	db         *db.Db
}

func NewServer(config config.Config, db *db.Db) (*Server, error) {
	timescaleRepo := timescale_repository.New(db.TimescalePool)
	postgresRepo := repository.New(db.PostgresPool)
	price.NewPriceService(timescaleRepo, postgresRepo)

	r := chi.NewRouter()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: r,
	}

	return &Server{
		httpServer: httpServer,
		router:     r,
		db:         db,
	}, nil
}

func (s *Server) Start() error {
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
