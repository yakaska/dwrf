package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Database struct {
	client *pgxpool.Pool
}

func Connect(ctx context.Context, dsn string) (*Database, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &Database{client: db}, nil
}

func (d *Database) Client() *pgxpool.Pool {
	return d.client
}

func (d *Database) Close() {
	d.client.Close()
}
