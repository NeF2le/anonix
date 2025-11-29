package grpc

import (
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateGRPC(grpcServer tokenizer.TokenizerServer) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryMiddleware(),
			interceptors.AddLogMiddleware,
		),
	)
	tokenizer.RegisterTokenizerServer(server, grpcServer)
	return server, nil
}

func CreateGRPCTLS(grpcServer tokenizer.TokenizerServer,
	tls credentials.TransportCredentials) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryMiddleware(),
			interceptors.AddLogMiddleware,
		),
		grpc.Creds(tls),
	)
	tokenizer.RegisterTokenizerServer(server, grpcServer)
	return server, nil
}
