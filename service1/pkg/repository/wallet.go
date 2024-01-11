package repository

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"service1/pkg/types"
)

func (p *Postgres) GetWalletBy(ctx context.Context, where map[string]any) (*types.Wallet, error) {
	var walletDB types.Wallet
	query, args, err := sq.Select(`id,actual,frozen`).
		From("wallet").
		Where(where).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Println("та самая")
		return nil, fmt.Errorf("unable to build SELECT query: %w", err)
	}
	err = p.conn.QueryRow(ctx, query, args...).Scan(&walletDB.ID, &walletDB.Actual, &walletDB.Frozen)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &walletDB, nil
}
