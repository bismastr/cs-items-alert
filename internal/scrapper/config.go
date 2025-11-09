package scrapper

import "time"

type Config struct {
	PageSize    int           `yaml:"page_size" env:"SCRAPER_PAGE_SIZE"`
	BaseUrl     string        `yaml:"base_url" env:"SCRAPER_BASE_URL"`
	TotalCount  int           `yaml:"total_count" env:"SCRAPER_TOTAL_COUNT"`
	BaseDelay   time.Duration `yaml:"base_delay" env:"SCRAPER_BASE_DELAY"`
	RandomDelay time.Duration `yaml:"random_delay" env:"SCRAPER_RANDOM_DELAY"`
	MaxRetries  int           `yaml:"max_retries" env:"SCRAPER_MAX_RETRIES"`
	BatchSize   int           `yaml:"batch_size" env:"SCRAPER_BATCH_SIZE"`
	MaxRateHits int           `yaml:"max_rate_hits" env:"SCRAPER_MAX_RATE_HITS"`
	MaxDelay    time.Duration `yaml:"max_delay" env:"SCRAPER_MAX_DELAY"`
	UserAgent   string        `yaml:"user_agent" env:"SCRAPER_USER_AGENT"`
}

func DefaultConfig() Config {
	return Config{
		PageSize:    10,
		BaseUrl:     "https://steamcommunity.com/market/search/render/?count=10&search_descriptions=0&sort_column=popular&sort_dir=desc&norender=1&category_730_Type=tag_CSGO_Type_WeaponCase&category_730_ItemSet%5B%5D=any&category_730_ProPlayer%5B%5D=any&category_730_StickerCapsule%5B%5D=any&category_730_Tournament%5B%5D=any&category_730_TournamentTeam%5B%5D=any&category_730_Type%5B%5D=tag_CSGO_Type_WeaponCase&category_730_Weapon%5B%5D=any&appid=730",
		TotalCount:  437,
		BaseDelay:   4 * time.Second,
		RandomDelay: 2 * time.Second,
		MaxRetries:  5,
		MaxRateHits: 3,
		UserAgent:   "Mozilla/5.0 (compatible; SteamScraper/1.0)",
	}
}
