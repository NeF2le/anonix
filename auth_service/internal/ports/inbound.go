package ports

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, login string, pass string, roleId int) (string, error)
	Login(ctx context.Context, login string, pass string) (string, string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)
}
