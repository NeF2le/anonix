package services

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/callers"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/gateway/internal/ports"
	"time"
)

type AuthService struct {
	AuthServiceRepo ports.AuthServiceRepository
	BaseDelay       time.Duration
	MaxRetries      uint
}

func NewAuthService(
	authServiceRepo ports.AuthServiceRepository,
	maxRetries uint,
	baseDelay time.Duration) *AuthService {
	return &AuthService{
		AuthServiceRepo: authServiceRepo,
		BaseDelay:       baseDelay,
		MaxRetries:      maxRetries,
	}
}

func (a *AuthService) Register(ctx context.Context, req *auth_service.RegisterRequest) (*auth_service.RegisterResponse, error) {
	resultChan := make(chan *auth_service.RegisterResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.Register(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call Register: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error) {
	resultChan := make(chan *auth_service.LoginResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.Login(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call Login: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) Refresh(ctx context.Context, req *auth_service.RefreshRequest) (*auth_service.RefreshResponse, error) {
	resultChan := make(chan *auth_service.RefreshResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.Refresh(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call Refresh: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, req *auth_service.IsAdminRequest) (*auth_service.IsAdminResponse, error) {
	resultChan := make(chan *auth_service.IsAdminResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.IsAdmin(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call IsAdmin: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) GetUserRoles(ctx context.Context, req *auth_service.GetUserRolesRequest) (*auth_service.GetUserRolesResponse, error) {
	resultChan := make(chan *auth_service.GetUserRolesResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.GetUserRoles(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetUserRoles: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) GetUsers(ctx context.Context, req *auth_service.GetUsersRequest) (*auth_service.GetUsersResponse, error) {
	resultChan := make(chan *auth_service.GetUsersResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.GetUsers(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetUsers: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*auth_service.DeleteUserResponse, error) {
	resultChan := make(chan *auth_service.DeleteUserResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.DeleteUser(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call DeleteUser: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) AssignRole(ctx context.Context, req *auth_service.AssignRoleRequest) (*auth_service.AssignRoleResponse, error) {
	resultChan := make(chan *auth_service.AssignRoleResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.AssignRole(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call AssignRole: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) RemoveRole(ctx context.Context, req *auth_service.RemoveRoleRequest) (*auth_service.RemoveRoleResponse, error) {
	resultChan := make(chan *auth_service.RemoveRoleResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.RemoveRole(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call RemoveRole: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) UpdateClearanceLevel(ctx context.Context, req *auth_service.UpdateClearanceLevelRequest) (*auth_service.UpdateClearanceLevelResponse, error) {
	resultChan := make(chan *auth_service.UpdateClearanceLevelResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.UpdateClearanceLevel(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateClearanceLevel: %w", err)
	}

	return <-resultChan, nil
}

func (a *AuthService) GetRolesList(ctx context.Context, req *auth_service.GetRolesListRequest) (*auth_service.GetRolesListResponse, error) {
	resultChan := make(chan *auth_service.GetRolesListResponse, 1)

	err := callers.Retry(func() error {
		resp, err := a.AuthServiceRepo.GetRolesList(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, a.MaxRetries, a.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetRolesList: %w", err)
	}

	return <-resultChan, nil
}
