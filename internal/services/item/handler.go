package item

import (
	"net/http"
	"strconv"

	"github.com/bismastr/cs-price-alert/internal/response"
	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	ItemService *ItemService
}

func NewItemHandler(itemService *ItemService) *ItemHandler {
	return &ItemHandler{
		ItemService: itemService,
	}
}

func (h *ItemHandler) GetItemDetails(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "item_id")
	if itemIDStr == "" {
		response.Error(w, 400, "1004", "Invalid item ID", map[string]interface{}{
			"error": "item_id is required",
		})
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 32)
	if err != nil {
		response.Error(w, 400, "1004", "Invalid item ID", map[string]interface{}{
			"error": "item_id must be a valid integer",
		})
		return
	}

	itemDetails, err := h.ItemService.GetItemDetails(r.Context(), int32(itemID))
	if err != nil {
		response.Error(w, 500, "1005", "Failed to get item details", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if itemDetails == nil {
		response.Error(w, 404, "1006", "Item not found", nil)
		return
	}

	response.Success(w, itemDetails)
}
