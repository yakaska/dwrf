package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"tinyurl/internal/model"
)

type Postgres struct {
	db *pgxpool.Pool
}

func NewPostgresDB(client *pgxpool.Pool) *Postgres {
	return &Postgres{db: client}
}

func (p *Postgres) Save(ctx context.Context, link model.Link) (*model.Link, error) {
	const op = "storage.postgres.Save"

	link.CreatedAt = time.Now().UTC()

	_, err := p.db.Exec(
		ctx,
		`INSERT INTO links (short, long, visits, created_at, expires_at, last_visited) VALUES ($1, $2, $3, $4, $5, $6)`,
		link.Short,
		link.Long,
		link.Visits,
		link.CreatedAt,
		link.ExpiresAt,
		link.LastVisited,
	)
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}

func (p *Postgres) Load(ctx context.Context, linkId string) (*model.Link, error) {
	const op = "storage.postgres.Load"

	var link model.Link

	err := p.db.QueryRow(
		ctx,
		`SELECT short, long, visits, created_at, expires_at, last_visited FROM links WHERE links.short = $1`,
		linkId,
	).Scan(
		&link.Short,
		&link.Long,
		&link.Visits,
		&link.CreatedAt,
		&link.ExpiresAt,
		&link.LastVisited,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, err
}

func (p *Postgres) AddVisits(ctx context.Context, linkId string) error {
	const op = "storage.postgres.AddVisits"

	_, err := p.db.Exec(
		ctx,
		`UPDATE links SET last_visited = NOW(), visits = visits + 1 WHERE links.short = $1`,
		linkId,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
