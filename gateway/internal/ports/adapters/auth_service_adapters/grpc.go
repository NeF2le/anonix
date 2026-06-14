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

func (a *AuthServiceAdapterGRPC) GetUsers(ctx context.Context, req *auth_service.GetUsersRequest) (*auth_service.GetUsersResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.GetUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*auth_service.DeleteUserResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.DeleteUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) AssignRole(ctx context.Context, req *auth_service.AssignRoleRequest) (*auth_service.AssignRoleResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.AssignRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) RemoveRole(ctx context.Context, req *auth_service.RemoveRoleRequest) (*auth_service.RemoveRoleResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.RemoveRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to remove role: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) UpdateClearanceLevel(ctx context.Context, req *auth_service.UpdateClearanceLevelRequest) (*auth_service.UpdateClearanceLevelResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.UpdateClearanceLevel(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update clearance level: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) GetRolesList(ctx context.Context, req *auth_service.GetRolesListRequest) (*auth_service.GetRolesListResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.GetRolesList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles list: %w", err)
	}
	return resp, nil
}

func (a *AuthServiceAdapterGRPC) GetUserRoles(ctx context.Context, req *auth_service.GetUserRolesRequest) (*auth_service.GetUserRolesResponse, error) {
	conn, err := grpc.NewClient(a.address, a.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for auth service: %w", err)
	}
	defer conn.Close()

	client := auth_service.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()

	resp, err := client.GetUserRoles(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	return resp, nil
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
