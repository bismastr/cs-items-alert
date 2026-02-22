package item

import "context"

type ItemDetailsResult struct {
	ItemId      int32  `json:"item_id"`
	Name        string `json:"name"`
	HashName    string `json:"hash_name"`
	IconUrl     string `json:"icon_url"`
	Price       int32  `json:"price"`
	SellListing int32  `json:"sell_listing"`
}

func (s *ItemService) GetItemDetails(ctx context.Context, itemId int32) (*ItemDetailsResult, error) {
	item, err := s.postgresRepo.GetItemByID(ctx, []int32{itemId})
	if err != nil {
		return nil, err
	}

	if len(item) == 0 {
		return nil, nil
	}

	price, err := s.timescaleRepo.GetItemLatestPrice(ctx, itemId)
	if err != nil {
		return nil, err
	}

	itemResult := item[0]
	return &ItemDetailsResult{
		ItemId:      itemResult.ID,
		Name:        itemResult.Name,
		HashName:    itemResult.HashName,
		IconUrl:     itemResult.IconUrl.String,
		Price:       price.SellPrice,
		SellListing: price.SellListings,
	}, nil
}
