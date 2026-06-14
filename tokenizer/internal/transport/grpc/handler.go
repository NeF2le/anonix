package grpc

import (
	"context"
	"errors"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/NeF2le/anonix/mapping/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type grpcTokenizerHandler struct {
	tokenizerClient ports.TokenizerUseCase
	tokenizer.UnimplementedTokenizerServer
}

func NewGRPCTokenizerHandler(tokenizerClient ports.TokenizerUseCase) tokenizer.TokenizerServer {
	return &grpcTokenizerHandler{tokenizerClient: tokenizerClient}
}

func (g *grpcTokenizerHandler) Tokenize(ctx context.Context, req *tokenizer.TokenizeRequest) (
	*tokenizer.TokenizeResponse, error) {
	if req.GetPlaintext() == nil {
		return nil, status.Error(codes.InvalidArgument, "plaintext is required")
	}

	pars := &domain.TokenizeParams{
		Plaintext:     req.GetPlaintext(),
		Deterministic: req.GetDeterministic(),
		Pseudonymize:  req.GetPseudonymize(),
		Algorithm:     req.GetAlgorithm(),
	}
	res, err := g.tokenizerClient.Tokenize(ctx, pars)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to tokenize",
			slog.Bool("deterministic", req.GetDeterministic()),
			slog.Bool("pseudonymize", req.GetPseudonymize()),
			logger.Err(err))
		return nil, status.Error(codes.Internal, "unable to tokenize plaintext")
	}
	logger.GetLoggerFromCtx(ctx).Debug(ctx,
		"tokenize result",
		slog.Bool("deterministic", req.GetDeterministic()),
		slog.String("algo", res.AlgoName))
	return &tokenizer.TokenizeResponse{
		TokenSuffix:   res.TokenSuffix,
		DekWrapped:    res.DekWrapped,
		CipherText:    res.Ciphertext,
		Deterministic: req.GetDeterministic(),
		AlgoName:      res.AlgoName,
	}, nil
}

func (g *grpcTokenizerHandler) Detokenize(ctx context.Context, req *tokenizer.DetokenizeRequest) (
	*tokenizer.DetokenizeResponse, error) {
	if req.GetCipherText() == nil {
		return nil, status.Error(codes.InvalidArgument, "ciphertext is required")
	}
	if req.GetDekWrapped() == nil {
		return nil, status.Error(codes.InvalidArgument, "dek wrapping is required")
	}

	pars := &domain.DetokenizeParams{
		Deterministic: req.GetDeterministic(),
		Ciphertext:    req.GetCipherText(),
		WrappedDek:    req.GetDekWrapped(),
		AlgoName:      req.GetAlgoName(),
	}

	plaintext, err := g.tokenizerClient.Detokenize(ctx, pars)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to detokenize",
			slog.Bool("deterministic", req.GetDeterministic()),
			logger.Err(err))
		if errors.Is(err, errs.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}
		return nil, status.Error(codes.Internal, "unable to detokenize token")
	}
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "detokenize result",
		slog.Bool("deterministic", req.GetDeterministic()),
		slog.Int("plaintext_len", len(plaintext)))

	return &tokenizer.DetokenizeResponse{Plaintext: plaintext}, nil
}

func (g *grpcTokenizerHandler) RotateMasterKey(ctx context.Context, _ *tokenizer.RotateMasterKeyRequest) (
	*tokenizer.RotateMasterKeyResponse, error) {
	if err := g.tokenizerClient.RotateMasterKey(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to rotate master key",
			logger.Err(err))
		return nil, status.Error(codes.Internal, "unable to rotate master key")
	}
	return &tokenizer.RotateMasterKeyResponse{}, nil
}

func (g *grpcTokenizerHandler) RewrapDEK(ctx context.Context, req *tokenizer.RewrapDEKRequest) (
	*tokenizer.RewrapDEKResponse, error) {
	if req.GetDekWrapped() == nil {
		return nil, status.Error(codes.InvalidArgument, "dek wrapping is required")
	}

	newWrappedDek, err := g.tokenizerClient.RewrapDEK(ctx, req.GetDekWrapped())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to rewrap dek",
			logger.Err(err))
		return nil, status.Error(codes.Internal, "unable to rewrap dek")
	}

	return &tokenizer.RewrapDEKResponse{DekWrapped: newWrappedDek}, nil
}

func (g *grpcTokenizerHandler) RotateDEK(ctx context.Context, req *tokenizer.RotateDEKRequest) (
	*tokenizer.RotateDEKResponse, error) {
	if req.GetDekWrapped() == nil {
		return nil, status.Error(codes.InvalidArgument, "dek wrapping is required")
	}
	if req.GetCipherText() == nil {
		return nil, status.Error(codes.InvalidArgument, "ciphertext is required")
	}

	pars := &domain.RotateDEKParams{
		WrappedDek:    req.GetDekWrapped(),
		Ciphertext:    req.GetCipherText(),
		Deterministic: req.GetDeterministic(),
		AlgoName:      req.GetAlgoName(),
	}

	res, err := g.tokenizerClient.RotateDEK(ctx, pars)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to rotate dek",
			slog.Bool("deterministic", req.GetDeterministic()),
			slog.String("algo_name", req.GetAlgoName()),
			logger.Err(err))
		return nil, status.Error(codes.Internal, "unable to rotate dek")
	}

	return &tokenizer.RotateDEKResponse{
		DekWrapped: res.DekWrapped,
		CipherText: res.Ciphertext,
		AlgoName:   res.AlgoName,
	}, nil
}
