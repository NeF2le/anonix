package runner

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"google.golang.org/grpc"
	"net"
)

func RunGRPC(ctx context.Context, grpcServer *grpc.Server, port int, serviceName string) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			fmt.Sprintf("%s: failed to listen", serviceName),
			logger.Err(err),
		)
		return fmt.Errorf("RunGRPC: failed to listen: %w", err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%s: listening on port %d", serviceName, port))
	if err = grpcServer.Serve(l); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			fmt.Sprintf("%s: failed to serve", serviceName),
			logger.Err(err),
		)
		return fmt.Errorf("RunGRPC: failed to serve: %w", err)
	}
	return nil
}

func MustRunGRPC(ctx context.Context, grpcServer *grpc.Server, port int, serviceName string) {
	if err := RunGRPC(ctx, grpcServer, port, serviceName); err != nil {
		panic(err)
	}
}
