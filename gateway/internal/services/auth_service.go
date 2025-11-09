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
