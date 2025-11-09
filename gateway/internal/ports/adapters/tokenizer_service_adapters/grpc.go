package tokenizer_service_adapters

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type TokenizerServiceAdapterGRPC struct {
	address     string
	opts        []grpc.DialOption
	dialTimeout time.Duration
}

func NewTokenizerServiceAdapterGRPC(address string, timeout time.Duration) *TokenizerServiceAdapterGRPC {
	return &TokenizerServiceAdapterGRPC{
		address:     address,
		dialTimeout: timeout,
		opts:        []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}
}

func (t *TokenizerServiceAdapterGRPC) Tokenize(ctx context.Context, req *tokenizer.TokenizeRequest) (
	*tokenizer.TokenizeResponse, error) {
	conn, err := grpc.NewClient(t.address, t.opts...)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to connect to tokenizer service",
			slog.String("address", t.address),
			slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to create new gRPC connection for tokenizer service: %w", err)
	}
	defer conn.Close()

	dctx, cancel := context.WithTimeout(ctx, t.dialTimeout)
	defer cancel()

	client := tokenizer.NewTokenizerClient(conn)
	resp, err := client.Tokenize(dctx, req)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to tokenize token",
			slog.String("address", t.address),
			slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to send gRPC request to tokenizer service: %w", err)
	}
	return resp, nil
}

func (t *TokenizerServiceAdapterGRPC) Detokenize(ctx context.Context, req *tokenizer.DetokenizeRequest) (
	*tokenizer.DetokenizeResponse, error) {
	conn, err := grpc.NewClient(t.address, t.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new gRPC connection for tokenizer service: %w", err)
	}
	defer conn.Close()

	dctx, cancel := context.WithTimeout(ctx, t.dialTimeout)
	defer cancel()

	client := tokenizer.NewTokenizerClient(conn)
	resp, err := client.Detokenize(dctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send gRPC request to tokenizer service: %w", err)
	}

	return resp, nil
}
