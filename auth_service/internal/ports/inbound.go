package ports

import (
	"context"
	"github.com/NeF2le/anonix/auth_service/internal/domain"
)

type AuthUseCase interface {
	Register(ctx context.Context, login string, pass string) (string, error)
	Login(ctx context.Context, login string, pass string) (string, string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)

	GetUsers(ctx context.Context) ([]*domain.User, error)
	DeleteUser(ctx context.Context, userId string) error

	AssignRole(ctx context.Context, userId string, roleId int) error
	RemoveRole(ctx context.Context, userId string, roleId int) error
	GetRolesList(ctx context.Context) ([]*domain.Role, error)
	GetUserRoles(ctx context.Context, userId string) ([]*domain.Role, error)

	UpdateClearanceLevel(ctx context.Context, userId string, level int) error
}
