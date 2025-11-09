package grpc

import (
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
)

func CreateGRPC(grpcServer auth_service.AuthServiceServer) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
	auth_service.RegisterAuthServiceServer(server, grpcServer)
	return server, nil
}
