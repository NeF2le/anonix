package grpc

import (
	"context"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/grpc/interceptors"
	"google.golang.org/grpc"
)

func CreateGRPC(ctx context.Context, grpcServer tokenizer.TokenizerServer) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryMiddleware(),
			interceptors.AddLogMiddleware,
		),
	)
	tokenizer.RegisterTokenizerServer(server, grpcServer)
	return server, nil
}
