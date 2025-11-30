package app

import (
	"log"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/server"
)

type ApiApp struct {
	Server *server.Server
}

func NewApiApp() (*ApiApp, error) {
	cfg := config.Load()

	db, err := db.NewDbClient(cfg)
	if err != nil {
		return nil, err
	}

	server, err := server.NewServer(cfg, db)
	if err != nil {
		return nil, err
	}

	return &ApiApp{
		Server: server,
	}, nil
}

func (app *ApiApp) Start() error {
	log.Printf("Starting server")
	return app.Server.Start()
}

func (app *ApiApp) Close() error {
	return app.Server.Close()
}
