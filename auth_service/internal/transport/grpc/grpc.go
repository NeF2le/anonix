package grpc

import (
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateGRPC(grpcServer auth_service.AuthServiceServer) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.RecoveryMiddleware(),
		interceptors.AddLogMiddleware))
	auth_service.RegisterAuthServiceServer(server, grpcServer)
	return server, nil
}

func CreateGRPCTLS(grpcServer auth_service.AuthServiceServer,
	tls credentials.TransportCredentials) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.RecoveryMiddleware(),
		interceptors.AddLogMiddleware),
		grpc.Creds(tls))
	auth_service.RegisterAuthServiceServer(server, grpcServer)
	return server, nil
}
