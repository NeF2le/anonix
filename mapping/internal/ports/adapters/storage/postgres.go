package storage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func NewPostgresAdapter(pool *pgxpool.Pool) *PostgresAdapter {
	return &PostgresAdapter{pool: pool}
}

func (p *PostgresAdapter) baseSelectMappingReq() sq.SelectBuilder {
	return sq.
		Select(
			"id",
			"dek_wrapped",
			"cipher_text",
			"token_ttl",
			"created_at",
			"deterministic",
			"reversible",
		).
		From("mapping.mappings").
		PlaceholderFormat(sq.Dollar)
}

func (p *PostgresAdapter) InsertMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error) {
	sql, args, err := sq.
		Insert("mapping.mappings").
		Columns("cipher_text", "dek_wrapped", "deterministic", "reversible", "token_ttl").
		Values(
			mapping.CipherText,
			mapping.DekWrapped,
			mapping.Deterministic,
			mapping.Reversible,
			mapping.TokenTtl,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InsertMapping: failed to build sql: %v", err)
	}

	err = p.pool.QueryRow(ctx, sql, args...).Scan(&mapping.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, errs.ErrMappingAlreadyExists
			}
		}
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("SQL: %s\nARGS: %v\n", sql, args))
		return nil, fmt.Errorf("InsertMapping: failed to scan id:  %v", err)
	}
	return mapping, nil
}

func (p *PostgresAdapter) DeleteMappingById(ctx context.Context, id uuid.UUID) error {
	_, err := p.SelectMappingById(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return errs.ErrMappingNotFound
	}

	sql, args, err := sq.
		Delete("mapping.mappings").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteMappingById: failed to build sql: %v", err)
	}
	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteMappingById: failed to execute sql: %v", err)
	}
	return nil
}

func (p *PostgresAdapter) SelectMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error) {
	sql, args, err := p.baseSelectMappingReq().Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetMappingById: failed to build sql: %v", err)
	}

	var mapping domain.Mapping
	err = p.pool.QueryRow(ctx, sql, args...).Scan(
		&mapping.ID,
		&mapping.DekWrapped,
		&mapping.CipherText,
		&mapping.TokenTtl,
		&mapping.CreatedAt,
		&mapping.Deterministic,
		&mapping.Reversible,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrMappingNotFound
		}
		return nil, fmt.Errorf("GetMappingById: failed to scan mapping: %v", err)
	}

	return &mapping, nil
}

func (p *PostgresAdapter) SelectAllMappings(ctx context.Context) ([]*domain.Mapping, error) {
	sql, args, err := p.baseSelectMappingReq().ToSql()
	if err != nil {
		return nil, fmt.Errorf("SelectAllMappings: failed to build sql: %v", err)
	}

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SelectAllMappings: failed to execute sql: %v", err)
	}
	defer rows.Close()

	var mappings []*domain.Mapping
	for rows.Next() {
		var mapping domain.Mapping
		err = rows.Scan(
			&mapping.ID,
			&mapping.DekWrapped,
			&mapping.CipherText,
			&mapping.TokenTtl,
			&mapping.CreatedAt,
			&mapping.Deterministic,
			&mapping.Reversible,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errs.ErrMappingNotFound
			}
			return nil, fmt.Errorf("SelectAllMappings: failed to scan mappings: %v", err)
		}
		mappings = append(mappings, &mapping)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("SelectAllMappings: rows iteration error: %v", err)
	}

	if len(mappings) == 0 {
		return nil, nil
	}

	return mappings, nil
}

func (p *PostgresAdapter) UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error) {
	sql, args, err := sq.
		Update("mapping.mappings").
		Set("token_ttl", tokenTtl).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UpdateMapping: failed to build sql: %v", err)
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UpdateMapping: failed to execute sql: %v", err)
	}

	mapping, err := p.SelectMappingById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UpdateMapping: failed to select updated mapping: %v", err)
	}

	return mapping, nil
}
