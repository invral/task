package common

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DependencyContainer struct {
	Pool   *pgxpool.Pool
	Config *Config
}

func NewDIContainer() (*DependencyContainer, error) {
	cfg := MustLoad()

	pool, err := NewConnectionDB(cfg.URL())
	if err != nil {
		return nil, err
	}

	return &DependencyContainer{
		Pool:   pool,
		Config: cfg,
	}, nil
}

func NewConnectionDB(uri string) (*pgxpool.Pool, error) {

	const op = "common.NewConnectionDB"

	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pool, nil
}
