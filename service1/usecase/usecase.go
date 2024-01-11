package usecase

import (
	"context"

	"service1/pkg/types"
)

type Repository interface {
	GetWalletBy(ctx context.Context, where map[string]any) (*types.Wallet, error)
	UpdateWallet(ctx context.Context, actual, frozen int, where map[string]any) error
	CreateTransaction(ctx context.Context, transactions *types.Transactions) (*string, *string, error)
	UpdateTransaction(ctx context.Context, status string, where map[string]any) error
}

type kafkaClient interface {
	SendMessage(msg []byte, topic string) error
}

type Config struct {
}

type Usecase struct {
	cfg        Config
	repository Repository
	kafka      kafkaClient
}

func New(cfg Config, database Repository, kafka kafkaClient) *Usecase {
	return &Usecase{
		cfg:        cfg,
		repository: database,
		kafka:      kafka,
	}
}
