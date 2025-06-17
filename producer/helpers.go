package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func OpenDB() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	return pool, nil
}
