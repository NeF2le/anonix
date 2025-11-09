package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/anonix/auth_service/internal/domain"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgresAdapter struct {
	pool *pgxpool.Pool
}

func NewAuthPostgresAdapter(pool *pgxpool.Pool) *AuthPostgresAdapter {
	return &AuthPostgresAdapter{pool: pool}
}

func (a *AuthPostgresAdapter) RegisterUser(ctx context.Context, login string, passHash []byte, roleId int) (string, error) {
	query := `INSERT INTO users.users (login, password_hash, role_id) VALUES ($1, $2, $3) RETURNING id`

	var userID string
	err := a.pool.QueryRow(ctx, query, login, passHash, roleId).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", errs.ErrUserAlreadyExists
		}
		return "", fmt.Errorf("error inserting user: %w", err)
	}

	return userID, nil
}

func (a *AuthPostgresAdapter) LoginUser(ctx context.Context, login string) (*domain.User, error) {
	query := `SELECT id, login, password_hash, role_id FROM users.users WHERE login = $1`

	var user domain.User
	err := a.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PassHash, &user.RoleId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &user, errs.ErrUserNotFound
		}
		return &user, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (a *AuthPostgresAdapter) IsAdminCheck(ctx context.Context, userId uuid.UUID) (bool, error) {
	query := `SELECT CASE WHEN role_id = 1 THEN TRUE ELSE FALSE END AS result FROM users.users WHERE id = $1`

	var result bool
	err := a.pool.QueryRow(ctx, query, userId).Scan(&result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errs.ErrUserNotFound
		}
		return false, fmt.Errorf("error checking user is admin: %w", err)
	}

	return result, nil
}
