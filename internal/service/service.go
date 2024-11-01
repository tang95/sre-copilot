package service

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type Service struct {
	config      *pkg.Config
	logger      *zap.Logger
	transaction Transaction
	cache       *cache.Cache
}

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}

func NewService(config *pkg.Config, logger *zap.Logger,
	transaction Transaction,
) (*Service, error) {
	return &Service{
		config:      config,
		logger:      logger,
		transaction: transaction,
		cache:       cache.New(cache.NoExpiration, cache.NoExpiration),
	}, nil
}
