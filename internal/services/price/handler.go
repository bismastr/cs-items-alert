package price

import (
	"net/http"

	"github.com/bismastr/cs-price-alert/internal/response"
)

type PriceHandler struct {
	priceService *PriceService
}

func NewPriceHandler(priceService *PriceService) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

func (h *PriceHandler) GetSearchPriceChanges(w http.ResponseWriter, r *http.Request) {

	searchResult, err := h.priceService.GetSearchPriceChanges(r.Context(), GetSearchPriceChanggesParams{
		Query: r.URL.Query().Get("query"),
		Limit: 10,
	})
	if err != nil {
		response.Error(w, 500, "1001", "Failed to get price changes", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.Success(w, searchResult)
}
