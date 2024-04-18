package repository

import (
	"javacode/internal/config"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type repImpl struct {
	cfg *config.Config
	log *slog.Logger
	DB  *sqlx.DB
}

func NewRep(cfg *config.Config, log *slog.Logger, dbx *sqlx.DB) *repImpl {
	return &repImpl{
		cfg: cfg,
		log: log,
		DB:  dbx,
	}
}

type Data struct {
	Role string
	Sum  int
}
