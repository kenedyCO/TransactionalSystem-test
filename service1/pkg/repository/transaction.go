package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"service1/pkg/types"
)

const (
	statusCreate = "Create"
)

func (p *Postgres) CreateTransaction(ctx context.Context, transactions *types.Transactions) (*string, *string, error) {
	query, args, err := sq.Insert("transactions").
		Columns(`amount,status`).
		Values(transactions.Amount, statusCreate).
		Suffix(`RETURNING id,status`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to build INSERT query: %w", err)
	}
	var id, status string
	row := p.conn.QueryRow(ctx, query, args...)
	err = row.Scan(&id, &status)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to scan: %w", err)
	}

	return &id, &status, nil
}

func (p *Postgres) UpdateTransaction(ctx context.Context, status string, where map[string]any) error {
	query, args, err := sq.Update("transactions").
		Set("status", status).
		Where(where).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("unable to build SELECT query: %w", err)
	}
	p.conn.QueryRow(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateWallet(ctx context.Context, actual, frozen int, where map[string]any) error {
	query, args, err := sq.Update("wallet").
		Set("actual", actual).
		Set("frozen", frozen).
		Where(where).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("unable to build SELECT query: %w", err)
	}
	p.conn.QueryRow(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
