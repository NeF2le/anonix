package mapping_service_adapters

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type MappingServiceAdapterGRPC struct {
	address     string
	opts        []grpc.DialOption
	dialTimeout time.Duration
}

func NewMappingServiceAdapterGRPC(address string, dialTimeout time.Duration) *MappingServiceAdapterGRPC {
	return &MappingServiceAdapterGRPC{
		address:     address,
		opts:        []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		dialTimeout: dialTimeout,
	}
}

func (s *MappingServiceAdapterGRPC) AddTLS(creds credentials.TransportCredentials) {
	s.opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
}

func (s *MappingServiceAdapterGRPC) CreateMapping(ctx context.Context, req *mapping.CreateMappingRequest) (
	*mapping.CreateMappingResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.CreateMapping(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check if token exists: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetMapping(ctx context.Context, req *mapping.GetMappingRequest) (
	*mapping.GetMappingResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetMapping(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetMappingByToken(ctx context.Context, req *mapping.GetMappingByTokenRequest) (
	*mapping.GetMappingResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetMappingByToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping by token: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) UpdateMapping(ctx context.Context, req *mapping.UpdateMappingRequest) (
	*mapping.UpdateMappingResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.UpdateMapping(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update mapping: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) UpdateMappingDek(ctx context.Context, req *mapping.UpdateMappingDekRequest) (
	*mapping.UpdateMappingDekResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.UpdateMappingDek(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update mapping dek: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) UpdateMappingCrypto(ctx context.Context, req *mapping.UpdateMappingCryptoRequest) (
	*mapping.UpdateMappingCryptoResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.UpdateMappingCrypto(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update mapping crypto: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) DeleteMapping(ctx context.Context, req *mapping.DeleteMappingRequest) (
	*mapping.DeleteMappingResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.DeleteMapping(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete mapping: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetMappingList(ctx context.Context, req *mapping.GetMappingListRequest) (
	*mapping.GetMappingListResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)
	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetMappingList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping list: %w", err)
	}
	return resp, nil
}

func (s *MappingServiceAdapterGRPC) CreateKind(ctx context.Context, req *mapping.CreateKindRequest) (
	*mapping.CreateKindResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.CreateKind(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create kind: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetKind(ctx context.Context, req *mapping.GetKindRequest) (
	*mapping.GetKindResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetKind(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get kind: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetKindByName(ctx context.Context, req *mapping.GetKindByNameRequest) (
	*mapping.GetKindByNameResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetKindByName(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get kind by name: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) ListKinds(ctx context.Context, req *mapping.ListKindsRequest) (
	*mapping.ListKindsResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.ListKinds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get kind list: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) UpdateKind(ctx context.Context, req *mapping.UpdateKindRequest) (
	*mapping.UpdateKindResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.UpdateKind(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update kind: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) DeleteKind(ctx context.Context, req *mapping.DeleteKindRequest) (
	*mapping.DeleteKindResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.DeleteKind(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete kind: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) CreateAuditLog(ctx context.Context, req *mapping.CreateAuditLogRequest) (
	*mapping.CreateAuditLogResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.CreateAuditLog(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit log entry: %w", err)
	}

	return resp, nil
}

func (s *MappingServiceAdapterGRPC) GetAuditLogList(ctx context.Context, req *mapping.GetAuditLogListRequest) (
	*mapping.GetAuditLogListResponse, error) {
	conn, err := grpc.NewClient(s.address, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for mapping service: %w", err)
	}
	defer conn.Close()

	client := mapping.NewMappingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, s.dialTimeout)
	defer cancel()

	resp, err := client.GetAuditLogList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log list: %w", err)
	}

	return resp, nil
}
