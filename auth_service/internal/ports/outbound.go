package ports

import (
	"context"
	"github.com/NeF2le/anonix/auth_service/internal/domain"
	"github.com/google/uuid"
)

type StorageRepository interface {
	RegisterUser(ctx context.Context, login string, passHash []byte) (string, error)
	LoginUser(ctx context.Context, login string) (*domain.User, error)
	IsAdminCheck(ctx context.Context, userId uuid.UUID) (bool, error)

	GetUsers(ctx context.Context) ([]*domain.User, error)
	DeleteUser(ctx context.Context, userId uuid.UUID) error

	AssignRole(ctx context.Context, userId uuid.UUID, roleId int) error
	RemoveRole(ctx context.Context, userId uuid.UUID, roleId int) error
	GetRolesList(ctx context.Context) ([]*domain.Role, error)
	GetUserRoles(ctx context.Context, userId uuid.UUID) ([]*domain.Role, error)

	UpdateClearanceLevel(ctx context.Context, userId uuid.UUID, level int) error
}

type CacheRepository interface {
	SaveToken(ctx context.Context, token, userID string, refresh bool) error
	GetToken(ctx context.Context, token string, refresh bool) (string, error)
	DeleteToken(ctx context.Context, token string, refresh bool) error
}
