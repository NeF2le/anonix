package auth_service_adapters

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthServiceAdapterGRPC struct {
	address     string
	opts        []grpc.DialOption
	dialTimeout time.Duration
}

func NewAuthServiceAdapterGRPC(address string, dialTimeout time.Duration) *AuthServiceAdapterGRPC {
	return &AuthServiceAdapterGRPC{
		address:     address,
		opts:        []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		dialTimeout: dialTimeout,
	}
}

func (a *AuthServiceAdapterGRPC) AddTLS(creds credentials.TransportCredentials) {
	a.opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
}

func (a *AuthServiceAdapterGRPC) Register(ctx context.Context, req *auth_service.RegisterRequest) (*auth_service.RegisterResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) Refresh(ctx context.Context, req *auth_service.RefreshRequest) (*auth_service.RefreshResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.Refresh(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) IsAdmin(ctx context.Context, req *auth_service.IsAdminRequest) (*auth_service.IsAdminResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.IsAdmin(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return resp, nil
}
