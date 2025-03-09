package steam

type SteamSearchResponse struct {
	Success    bool       `json:"success"`
	Start      int        `json:"start"`
	PageSize   int        `json:"pagesize"`
	TotalCount int        `json:"total_count"`
	SearchData SearchData `json:"searchdata"`
	Results    []Result   `json:"results"`
}

type SearchData struct {
	Query              string `json:"query"`
	SearchDescriptions bool   `json:"search_descriptions"`
	TotalCount         int    `json:"total_count"`
	PageSize           int    `json:"pagesize"`
	Prefix             string `json:"prefix"`
	ClassPrefix        string `json:"class_prefix"`
}

type Result struct {
	Name             string           `json:"name"`
	HashName         string           `json:"hash_name"`
	SellListings     int              `json:"sell_listings"`
	SellPrice        int              `json:"sell_price"`
	SellPriceText    string           `json:"sell_price_text"`
	AppIcon          string           `json:"app_icon"`
	AppName          string           `json:"app_name"`
	AssetDescription AssetDescription `json:"asset_description"`
	SalePriceText    string           `json:"sale_price_text"`
}

type AssetDescription struct {
	AppID           int    `json:"appid"`
	ClassID         string `json:"classid"`
	InstanceID      string `json:"instanceid"`
	BackgroundColor string `json:"background_color"`
	IconURL         string `json:"icon_url"`
	Tradable        int    `json:"tradable"`
	Name            string `json:"name"`
	NameColor       string `json:"name_color"`
	Type            string `json:"type"`
	MarketName      string `json:"market_name"`
	MarketHashName  string `json:"market_hash_name"`
	Commodity       int    `json:"commodity"`
}
