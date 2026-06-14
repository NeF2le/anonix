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
			"m.id",
			"m.token",
			"m.dek_wrapped",
			"m.cipher_text",
			"m.token_ttl",
			"m.created_at",
			"m.deterministic",
			"m.algo_name",
			"k.id AS kind_id",
			"k.name AS kind_name",
			"k.access_level",
			"k.russian_name",
			"k.short_name",
		).
		From("mapping.mappings m").
		LeftJoin("mapping.kinds k ON k.id = m.kind_id").
		PlaceholderFormat(sq.Dollar)
}

func (p *PostgresAdapter) baseSelectAuditLogReq() sq.SelectBuilder {
	return sq.
		Select(
			"a.id",
			"a.user_id",
			"a.action",
			"a.token",
			"a.created_at",
			"k.id AS kind_id",
			"k.name AS kind_name",
			"k.access_level",
			"k.russian_name",
			"k.short_name",
		).
		From("mapping.audit_log a").
		LeftJoin("mapping.kinds k ON k.id = a.kind_id").
		PlaceholderFormat(sq.Dollar)
}

func (p *PostgresAdapter) baseSelectKindReq() sq.SelectBuilder {
	return sq.
		Select(
			"id",
			"name",
			"russian_name",
			"access_level",
			"mask",
			"short_name",
		).
		From("mapping.kinds").
		PlaceholderFormat(sq.Dollar)
}

func (p *PostgresAdapter) InsertMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error) {
	var kindID *int32
	if mapping.Kind != nil {
		kindID = &mapping.Kind.Id
	}

	sql, args, err := sq.
		Insert("mapping.mappings").
		Columns("token", "cipher_text", "dek_wrapped", "deterministic", "kind_id", "token_ttl", "algo_name").
		Values(
			mapping.Token,
			mapping.CipherText,
			mapping.DekWrapped,
			mapping.Deterministic,
			kindID,
			mapping.TokenTtl,
			mapping.AlgoName,
		).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InsertMapping: failed to build sql: %v", err)
	}

	err = p.pool.QueryRow(ctx, sql, args...).Scan(&mapping.ID, &mapping.CreatedAt)
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
	sql, args, err := p.baseSelectMappingReq().Where(sq.Eq{"m.id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetMappingById: failed to build sql: %v", err)
	}

	var mapping domain.Mapping
	var (
		kindID      *int32
		kindName    *string
		accessLevel *int32
		russianName *string
		shortName   *string
	)

	err = p.pool.QueryRow(ctx, sql, args...).Scan(
		&mapping.ID,
		&mapping.Token,
		&mapping.DekWrapped,
		&mapping.CipherText,
		&mapping.TokenTtl,
		&mapping.CreatedAt,
		&mapping.Deterministic,
		&mapping.AlgoName,
		&kindID,
		&kindName,
		&accessLevel,
		&russianName,
		&shortName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrMappingNotFound
		}
		return nil, fmt.Errorf("GetMappingById: failed to scan mapping: %v", err)
	}
	if kindID != nil {
		mapping.Kind = &domain.Kind{
			Id:          *kindID,
			Name:        *kindName,
			AccessLevel: *accessLevel,
			RussianName: *russianName,
			ShortName:   *shortName,
		}
	}

	return &mapping, nil
}

func (p *PostgresAdapter) SelectMappingByToken(ctx context.Context, token string) (*domain.Mapping, error) {
	sql, args, err := p.baseSelectMappingReq().Where(sq.Eq{"m.token": token}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("SelectMappingByToken: failed to build sql: %v", err)
	}

	var mapping domain.Mapping
	var (
		kindID      *int32
		kindName    *string
		accessLevel *int32
		russianName *string
		shortName   *string
	)

	err = p.pool.QueryRow(ctx, sql, args...).Scan(
		&mapping.ID,
		&mapping.Token,
		&mapping.DekWrapped,
		&mapping.CipherText,
		&mapping.TokenTtl,
		&mapping.CreatedAt,
		&mapping.Deterministic,
		&mapping.AlgoName,
		&kindID,
		&kindName,
		&accessLevel,
		&russianName,
		&shortName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrMappingNotFound
		}
		return nil, fmt.Errorf("SelectMappingByToken: failed to scan mapping: %v", err)
	}
	if kindID != nil {
		mapping.Kind = &domain.Kind{
			Id:          *kindID,
			Name:        *kindName,
			AccessLevel: *accessLevel,
			RussianName: *russianName,
			ShortName:   *shortName,
		}
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
		var (
			kindID      *int32
			kindName    *string
			accessLevel *int32
			russianName *string
			shortName   *string
		)

		err = rows.Scan(
			&mapping.ID,
			&mapping.Token,
			&mapping.DekWrapped,
			&mapping.CipherText,
			&mapping.TokenTtl,
			&mapping.CreatedAt,
			&mapping.Deterministic,
			&mapping.AlgoName,
			&kindID,
			&kindName,
			&accessLevel,
			&russianName,
			&shortName,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errs.ErrMappingNotFound
			}
			return nil, fmt.Errorf("SelectAllMappings: failed to scan mappings: %v", err)
		}

		if kindID != nil {
			mapping.Kind = &domain.Kind{
				Id:          *kindID,
				Name:        *kindName,
				AccessLevel: *accessLevel,
				RussianName: *russianName,
				ShortName:   *shortName,
			}
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

func (p *PostgresAdapter) UpdateMappingDek(ctx context.Context, id uuid.UUID, dekWrapped []byte) error {
	sql, args, err := sq.
		Update("mapping.mappings").
		Set("dek_wrapped", dekWrapped).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("UpdateMappingDek: failed to build sql: %v", err)
	}

	if _, err = p.pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("UpdateMappingDek: failed to execute sql: %v", err)
	}

	return nil
}

func (p *PostgresAdapter) UpdateMappingCrypto(ctx context.Context, id uuid.UUID, dekWrapped, cipherText []byte, algoName string) error {
	sql, args, err := sq.
		Update("mapping.mappings").
		Set("dek_wrapped", dekWrapped).
		Set("cipher_text", cipherText).
		Set("algo_name", algoName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("UpdateMappingCrypto: failed to build sql: %v", err)
	}

	if _, err = p.pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("UpdateMappingCrypto: failed to execute sql: %v", err)
	}

	return nil
}

func (p *PostgresAdapter) InsertKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error) {
	sql, args, err := sq.
		Insert("mapping.kinds").
		Columns(
			"name",
			"russian_name",
			"access_level",
			"mask",
			"short_name",
		).
		Values(
			kind.Name,
			kind.RussianName,
			kind.AccessLevel,
			kind.Mask,
			kind.ShortName,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InsertKind: failed to build sql: %v", err)
	}

	err = p.pool.QueryRow(ctx, sql, args...).Scan(&kind.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, errs.ErrKindAlreadyExists
			}
		}

		logger.GetLoggerFromCtx(ctx).Info(
			ctx,
			fmt.Sprintf("SQL: %s\nARGS: %v\n", sql, args),
		)

		return nil, fmt.Errorf("InsertKind: failed to scan id: %v", err)
	}

	return kind, nil
}

func (p *PostgresAdapter) SelectKindById(ctx context.Context, id int32) (*domain.Kind, error) {
	sql, args, err := p.baseSelectKindReq().
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SelectKindById: failed to build sql: %v", err)
	}

	var kind domain.Kind

	err = p.pool.QueryRow(ctx, sql, args...).Scan(
		&kind.Id,
		&kind.Name,
		&kind.RussianName,
		&kind.AccessLevel,
		&kind.Mask,
		&kind.ShortName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrKindNotFound
		}

		return nil, fmt.Errorf("SelectKindById: failed to scan kind: %v", err)
	}

	return &kind, nil
}

func (p *PostgresAdapter) SelectKindByName(ctx context.Context, name string) (*domain.Kind, error) {
	sql, args, err := p.baseSelectKindReq().
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SelectKindByName: failed to build sql: %v", err)
	}

	var kind domain.Kind

	err = p.pool.QueryRow(ctx, sql, args...).Scan(
		&kind.Id,
		&kind.Name,
		&kind.RussianName,
		&kind.AccessLevel,
		&kind.Mask,
		&kind.ShortName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrKindNotFound
		}

		return nil, fmt.Errorf("SelectKindByName: failed to scan kind: %v", err)
	}

	return &kind, nil
}

func (p *PostgresAdapter) GetAllKinds(ctx context.Context) ([]*domain.Kind, error) {
	sql, args, err := p.baseSelectKindReq().ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetAllKinds: failed to build sql: %v", err)
	}

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAllKinds: failed to execute sql: %v", err)
	}
	defer rows.Close()

	var kinds []*domain.Kind

	for rows.Next() {
		var kind domain.Kind

		err = rows.Scan(
			&kind.Id,
			&kind.Name,
			&kind.RussianName,
			&kind.AccessLevel,
			&kind.Mask,
			&kind.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllKinds: failed to scan kind: %v", err)
		}

		kinds = append(kinds, &kind)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllKinds: rows iteration error: %v", err)
	}

	if len(kinds) == 0 {
		return nil, nil
	}

	return kinds, nil
}

func (p *PostgresAdapter) UpdateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error) {
	sql, args, err := sq.
		Update("mapping.kinds").
		Set("name", kind.Name).
		Set("russian_name", kind.RussianName).
		Set("access_level", kind.AccessLevel).
		Set("mask", kind.Mask).
		Set("short_name", kind.ShortName).
		Where(sq.Eq{"id": kind.Id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UpdateKind: failed to build sql: %v", err)
	}

	tag, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UpdateKind: failed to execute sql: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return nil, errs.ErrKindNotFound
	}

	return p.SelectKindById(ctx, kind.Id)
}

func (p *PostgresAdapter) DeleteKindById(ctx context.Context, id int32) error {
	_, err := p.SelectKindById(ctx, id)
	if err != nil {
		return err
	}

	sql, args, err := sq.
		Delete("mapping.kinds").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteKindById: failed to build sql: %v", err)
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return errs.ErrKindInUse
			}
		}
		return fmt.Errorf("DeleteKindById: failed to execute sql: %v", err)
	}

	return nil
}

func (p *PostgresAdapter) CreateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error) {
	return p.InsertKind(ctx, kind)
}

func (p *PostgresAdapter) GetKindById(ctx context.Context, id int32) (*domain.Kind, error) {
	return p.SelectKindById(ctx, id)
}

func (p *PostgresAdapter) GetKindByName(ctx context.Context, name string) (*domain.Kind, error) {
	return p.SelectKindByName(ctx, name)
}

func (p *PostgresAdapter) CreateAuditLog(ctx context.Context, entry *domain.AuditLogEntry) (*domain.AuditLogEntry, error) {
	var kindID *int32
	if entry.Kind != nil {
		kindID = &entry.Kind.Id
	}

	sql, args, err := sq.
		Insert("mapping.audit_log").
		Columns("user_id", "action", "token", "kind_id").
		Values(entry.UserID, entry.Action, entry.Token, kindID).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateAuditLog: failed to build sql: %v", err)
	}

	err = p.pool.QueryRow(ctx, sql, args...).Scan(&entry.ID, &entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateAuditLog: failed to scan id: %v", err)
	}

	return entry, nil
}

func (p *PostgresAdapter) GetAuditLogList(ctx context.Context) ([]*domain.AuditLogEntry, error) {
	sql, args, err := p.baseSelectAuditLogReq().
		OrderBy("a.created_at DESC").
		Limit(500).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetAuditLogList: failed to build sql: %v", err)
	}

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAuditLogList: failed to execute sql: %v", err)
	}
	defer rows.Close()

	var entries []*domain.AuditLogEntry
	for rows.Next() {
		var entry domain.AuditLogEntry
		var (
			kindID      *int32
			kindName    *string
			accessLevel *int32
			russianName *string
			shortName   *string
		)

		err = rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Action,
			&entry.Token,
			&entry.CreatedAt,
			&kindID,
			&kindName,
			&accessLevel,
			&russianName,
			&shortName,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAuditLogList: failed to scan entry: %v", err)
		}

		if kindID != nil {
			entry.Kind = &domain.Kind{
				Id:          *kindID,
				Name:        *kindName,
				AccessLevel: *accessLevel,
				RussianName: *russianName,
				ShortName:   *shortName,
			}
		}

		entries = append(entries, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAuditLogList: rows iteration error: %v", err)
	}

	return entries, nil
}
