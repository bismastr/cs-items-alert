package price

type PriceHandler struct {
	priceService *PriceService
}

func NewPriceHandler(priceService *PriceService) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}
