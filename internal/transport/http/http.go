package http

import (
	"context"
	"javacode/internal/config"
	"javacode/internal/service"
	"log/slog"
	"net/http"
)

type Service interface {
	ChangingSum(ctx context.Context, data *service.TrasactionInfo) (int, error)
}

type ApiImpl struct {
	cfg     *config.Config
	log     *slog.Logger
	service Service
	server  *http.Server
}

// NewAPI
func NewApi(cfg *config.Config, log *slog.Logger, service Service) *ApiImpl {
	return &ApiImpl{
		service: service,
		log:     log,
		cfg:     cfg,
	}
}
