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
		SortBy: r.URL.Query().Get("sort_by"),
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

func (h *PriceHandler) GetItemPriceChart(w http.ResponseWriter, r *http.Request) {
	itemID := utils.GetQueryInt(r, "item_id", 0)
	if itemID == 0 {
		response.Error(w, 400, "1002", "Invalid item ID", map[string]interface{}{
			"error": "item_id is required",
		})
		return
	}
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		response.Error(w, 400, "1003", "Invalid interval", map[string]interface{}{
			"error": "interval is required",
		})
		return
	}

	chartData, err := h.priceService.GetItemPriceChart(r.Context(), itemID, ChartPeriod(interval))
	if err != nil {
		response.Error(w, 500, "1004", "Failed to get item price chart", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.Success(w, chartData)
}

func (h *PriceHandler) GetItemPriceStats(w http.ResponseWriter, r *http.Request) {
	itemID := utils.GetQueryInt(r, "item_id", 0)
	if itemID == 0 {
		response.Error(w, 400, "1005", "Invalid item ID", map[string]interface{}{
			"error": "item_id is required",
		})
		return
	}

	statsData, err := h.priceService.GetItemPriceStats(r.Context(), itemID)
	if err != nil {
		response.Error(w, 500, "1007", "Failed to get item price stats", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.Success(w, statsData)
}
