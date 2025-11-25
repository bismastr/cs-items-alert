package server

import (
	"github.com/bismastr/cs-price-alert/internal/services/price"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	PriceService *price.PriceService
}

func NewRouter(priceService *price.PriceService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	return r
}
