package service

import (
	"context"
	"javacode/internal/config"
	"log/slog"
)

type Repository interface {
	SumInc(ctx context.Context, sum int) (int, error)
	SumDec(ctx context.Context, sum int) (int, error)
}

type servImpl struct {
	cfg *config.Config
	log *slog.Logger
	rep Repository
}

func NewServ(
	cfg *config.Config,
	log *slog.Logger,
	rep Repository,
) *servImpl {
	return &servImpl{
		rep: rep,
		log: log,
		cfg: cfg,
	}
}

type TrasactionInfo struct {
	Role string `json:"role"`
	Sum  int    `json:"sum"`
}

type RespTrasactionInfo struct {
	Sum int `json:"sum"`
}
