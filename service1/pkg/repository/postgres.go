package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnString string
}

type Postgres struct {
	cfg  Config
	conn *pgxpool.Pool
}

func New(cfg Config) *Postgres {
	return &Postgres{
		cfg: cfg,
	}
}

func (p *Postgres) Start(ctx context.Context) error {
	dbpool, err := pgxpool.New(ctx, p.cfg.ConnString)
	if err != nil {
		return err
	}
	p.conn = dbpool
	log.Println("NewConnToDB end")

	return nil
}

func (p *Postgres) ShutDown(context.Context) error {
	p.conn.Close()
	log.Println("CloseConnToDB end")

	return nil
}
