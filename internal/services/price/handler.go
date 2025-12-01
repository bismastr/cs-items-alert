package price

import (
	"net/http"

	"github.com/bismastr/cs-price-alert/internal/response"
	"github.com/bismastr/cs-price-alert/internal/utils"
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
	page := utils.GetQueryInt(r, "page", 1)
	limit := utils.GetQueryInt(r, "limit", 16)
	offset := (page - 1) * limit

	searchResult, totalCount, err := h.priceService.GetSearchPriceChanges(r.Context(), PriceChangeQueryParams{
		Query:  r.URL.Query().Get("query"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, 500, "1001", "Failed to get price changes", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.Success(w, map[string]interface{}{
		"total": totalCount,
		"items": searchResult,
	})
}
