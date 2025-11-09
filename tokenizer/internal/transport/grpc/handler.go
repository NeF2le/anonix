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
		Reversible:    req.GetReversible(),
	}
	res, err := g.tokenizerClient.Tokenize(ctx, pars)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to tokenize",
			slog.String("plaintext", string(req.GetPlaintext())),
			slog.Bool("deterministic", req.GetDeterministic()),
			slog.Bool("reversible", req.GetReversible()),
			logger.Err(err))
		return nil, status.Error(codes.Internal, "unable to tokenize plaintext")
	}
	logger.GetLoggerFromCtx(ctx).Debug(ctx,
		"tokenize result",
		slog.String("plaintext", string(req.GetPlaintext())),
		slog.Bool("deterministic", req.GetDeterministic()),
		slog.Bool("reversible", req.GetReversible()),
		logger.Base64("dek wrapped", res.DekWrapped),
		logger.Base64("ciphertext", res.Ciphertext))
	return &tokenizer.TokenizeResponse{
		DekWrapped:    res.DekWrapped,
		CipherText:    res.Ciphertext,
		Deterministic: req.GetDeterministic(),
		Reversible:    req.GetReversible(),
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
	}

	plaintext, err := g.tokenizerClient.Detokenize(ctx, pars)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to detokenize",
			logger.Base64("ciphertext", req.GetCipherText()),
			logger.Base64("dek wrapped", req.GetDekWrapped()),
			slog.Bool("deterministic", req.GetDeterministic()),
			logger.Err(err))
		if errors.Is(err, errs.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}
		return nil, status.Error(codes.Internal, "unable to detokenize token")
	}
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "detokenize result",
		slog.String("plaintext", string(plaintext)))

	return &tokenizer.DetokenizeResponse{Plaintext: plaintext}, nil
}
