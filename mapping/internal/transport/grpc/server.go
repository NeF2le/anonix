package grpc

import (
	"context"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
)

func CreateGRPC(ctx context.Context, grpcServer mapping.MappingServer) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryMiddleware(),
			interceptors.AddLogMiddleware,
		),
	)
	mapping.RegisterMappingServer(server, grpcServer)
	return server, nil
}
