package interceptors

import (
	"context"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

func AddLogMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx = logger.New(ctx)
	ctx = context.WithValue(ctx, logger.KeyForRequestID, uuid.New().String())
	logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC request",
		slog.String("method", info.FullMethod),
		slog.String("path", info.FullMethod),
		slog.Time("request time", time.Now()),
	)
	reply, err := handler(ctx, req)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "gRPC handler returns error", logger.Err(err))
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC reply")
	return reply, err
}

func RecoveryMiddleware() grpc.UnaryServerInterceptor {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) error {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	return recovery.UnaryServerInterceptor(recoveryOpts...)
}
