package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database    DatabaseConfig
	TimescaleDB DatabaseConfig
	Scraper     ScraperConfig
	RabbitMQ    RabbitMQConfig
	Server      ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	MaxConns int32
	MinConns int32
}

type ScraperConfig struct {
	PageSize    int
	BaseURL     string
	TotalCount  int
	BaseDelay   time.Duration
	RandomDelay time.Duration
	MaxRetries  int
}

type RabbitMQConfig struct {
	URL      string
	Host     string
	Username string
	Password string
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := Config{
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetString("database.port"),
			Username: v.GetString("database.username"),
			Password: v.GetString("database.password"),
			Database: v.GetString("database.name"),
			MaxConns: int32(v.GetInt("database.max_conns")),
			MinConns: int32(v.GetInt("database.min_conns")),
		},
		TimescaleDB: DatabaseConfig{
			Host:     v.GetString("timescaledb.host"),
			Port:     v.GetString("timescaledb.port"),
			Username: v.GetString("timescaledb.username"),
			Password: v.GetString("timescaledb.password"),
			Database: v.GetString("timescaledb.name"),
			MaxConns: int32(v.GetInt("timescaledb.max_conns")),
			MinConns: int32(v.GetInt("timescaledb.min_conns")),
		},
		Scraper: ScraperConfig{
			PageSize:    v.GetInt("scraper.pagesize"),
			BaseURL:     v.GetString("scraper.baseurl"),
			TotalCount:  v.GetInt("scraper.totalcount"),
			BaseDelay:   v.GetDuration("scraper.basedelay"),
			RandomDelay: v.GetDuration("scraper.randomdelay"),
			MaxRetries:  v.GetInt("scraper.maxretries"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      v.GetString("rmq.url"),
			Host:     v.GetString("rmq.host"),
			Username: v.GetString("rmq.username"),
			Password: v.GetString("rmq.password"),
		},
		Server: ServerConfig{
			Port: v.GetString("server.port"),
		},
	}

	return &cfg
}
