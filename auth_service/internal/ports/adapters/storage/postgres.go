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

func (a *AuthPostgresAdapter) GetUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.login, u.password_hash, u.clearance_level, ur.role_id, r.name AS role_name
		FROM auth.users u
		LEFT JOIN auth.users_roles ur ON ur.user_id = u.id
		LEFT JOIN auth.roles r ON r.id = ur.role_id
		ORDER BY u.login
	`

	rows, err := a.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	usersMap := make(map[string]*domain.User)
	for rows.Next() {
		var (
			userID         string
			login          string
			passHash       []byte
			clearanceLevel int
			roleID         *int
			roleName       *string
		)

		err = rows.Scan(&userID, &login, &passHash, &clearanceLevel, &roleID, &roleName)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		user, exists := usersMap[userID]
		if !exists {
			user = &domain.User{
				ID:             userID,
				Login:          login,
				PassHash:       passHash,
				Roles:          []*domain.Role{},
				ClearanceLevel: clearanceLevel,
			}
			usersMap[userID] = user
		}

		if roleID != nil && roleName != nil {
			user.Roles = append(user.Roles, &domain.Role{
				ID:   *roleID,
				Name: *roleName,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	result := make([]*domain.User, 0, len(usersMap))
	for _, u := range usersMap {
		result = append(result, u)
	}

	return result, nil
}

func (a *AuthPostgresAdapter) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	query := `DELETE FROM auth.users WHERE id = $1`

	cmdTag, err := a.pool.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (a *AuthPostgresAdapter) AssignRole(ctx context.Context, userId uuid.UUID, roleId int) error {
	query := `
		INSERT INTO auth.users_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := a.pool.Exec(ctx, query, userId, roleId)
	if err != nil {
		return fmt.Errorf("assign role: %w", err)
	}

	return nil
}

func (a *AuthPostgresAdapter) RemoveRole(ctx context.Context, userId uuid.UUID, roleId int) error {
	query := `
		DELETE FROM auth.users_roles
		WHERE user_id = $1 AND role_id = $2
	`

	cmdTag, err := a.pool.Exec(ctx, query, userId, roleId)
	if err != nil {
		return fmt.Errorf("remove role: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (a *AuthPostgresAdapter) UpdateClearanceLevel(ctx context.Context, userId uuid.UUID, level int) error {
	query := `UPDATE auth.users SET clearance_level = $1 WHERE id = $2`

	cmdTag, err := a.pool.Exec(ctx, query, level, userId)
	if err != nil {
		return fmt.Errorf("update clearance level: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (a *AuthPostgresAdapter) GetRolesList(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT id, name FROM auth.roles ORDER BY id`

	rows, err := a.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get roles list: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role

	for rows.Next() {
		var r domain.Role

		if err = rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}

		roles = append(roles, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return roles, nil
}

func NewAuthPostgresAdapter(pool *pgxpool.Pool) *AuthPostgresAdapter {
	return &AuthPostgresAdapter{pool: pool}
}

func (a *AuthPostgresAdapter) RegisterUser(ctx context.Context, login string, passHash []byte) (string, error) {
	query := `INSERT INTO auth.users (login, password_hash) VALUES ($1, $2) RETURNING id`

	var userID string
	err := a.pool.QueryRow(ctx, query, login, passHash).Scan(&userID)
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
	query := `SELECT id, login, password_hash, clearance_level FROM auth.users WHERE login = $1`

	var user domain.User
	err := a.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PassHash, &user.ClearanceLevel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &user, errs.ErrUserNotFound
		}
		return &user, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (a *AuthPostgresAdapter) IsAdminCheck(ctx context.Context, userId uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) AS result FROM auth.users_roles WHERE user_id = $1 AND role_id = 1`

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

func (a *AuthPostgresAdapter) GetUserRoles(ctx context.Context, userId uuid.UUID) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name
		FROM auth.users_roles ur
		INNER JOIN auth.roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := a.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("get roles list: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role

	for rows.Next() {
		var r domain.Role

		if err = rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}

		roles = append(roles, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return roles, nil
}
