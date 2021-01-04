package service

import (
	"context"
	"github.com/nokrPOC/internal/config"
)

type HermesService struct {}

func NewService(ctx context.Context, config *config.HermesServiceConfig) (*HermesService, error) {
	return &HermesService{}, nil
}