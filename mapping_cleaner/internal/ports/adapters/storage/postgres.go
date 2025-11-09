package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func NewPostgresAdapter(pool *pgxpool.Pool) *PostgresAdapter {
	return &PostgresAdapter{pool: pool}
}

func (p *PostgresAdapter) DeleteExpiredMappings(ctx context.Context) ([]uuid.UUID, error) {
	query := `
		DELETE FROM mapping.mappings
		WHERE (created_at + (token_ttl / 1000000000 * interval '1 second')) < now()
		RETURNING id`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to delete expired mappings from storage: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan id from storage: %w", err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return ids, nil
}
