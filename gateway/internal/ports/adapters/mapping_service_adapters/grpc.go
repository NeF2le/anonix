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
