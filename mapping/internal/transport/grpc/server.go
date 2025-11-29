package grpc

import (
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateGRPC(grpcServer mapping.MappingServer) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryMiddleware(),
			interceptors.AddLogMiddleware,
		),
	)
	mapping.RegisterMappingServer(server, grpcServer)
	return server, nil
}

func CreateGRPCTLS(grpcServer mapping.MappingServer,
	tls credentials.TransportCredentials) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.RecoveryMiddleware(),
		interceptors.AddLogMiddleware),
		grpc.Creds(tls))
	mapping.RegisterMappingServer(server, grpcServer)
	return server, nil
}
