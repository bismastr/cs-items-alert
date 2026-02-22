package server

import (
	"net/http"

	"github.com/bismastr/cs-price-alert/internal/response"
	"github.com/bismastr/cs-price-alert/internal/services/price"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(priceHandler *price.PriceHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	// Add Keep-Alive middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Keep-Alive", "timeout=60, max=100")
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/price-changes/search", priceHandler.GetSearchPriceChanges)
		r.Get("/price-changes/chart", priceHandler.GetItemPriceChart)
		r.Get("/price-changes/price-stats", priceHandler.GetItemPriceStats)
	})

	return r
}
