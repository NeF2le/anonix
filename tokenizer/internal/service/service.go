package service

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/NeF2le/anonix/mapping/internal/ports"
	"github.com/NeF2le/anonix/mapping/internal/service/algorithms"
	"log/slog"
)

type TokenizerService struct {
	vault         ports.VaultRepository
	convergentKey string
	dekBitsLength int
	jwtSecret     string
}

func NewTokenizerService(
	vault ports.VaultRepository,
	convergentKey string,
	dekBitsLength int) *TokenizerService {
	return &TokenizerService{
		vault:         vault,
		convergentKey: convergentKey,
		dekBitsLength: dekBitsLength,
	}
}

func (t *TokenizerService) Tokenize(ctx context.Context, pars *domain.TokenizeParams) (*domain.TokenResult, error) {
	deterministic, reversible := pars.Deterministic, pars.Reversible
	var err error
	var dek []byte
	var wrappedDek []byte
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(dek)

	if deterministic && reversible {
		wrappedDek, dek, err = t.vault.GenerateDEK(ctx, t.dekBitsLength, t.convergentKey, "test")
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				"failed to generate DEK",
				slog.Int("dek_bits_length", t.dekBitsLength),
				slog.String("key", t.convergentKey),
				logger.Err(err),
			)
			return nil, fmt.Errorf("failed to generate DEK: %w", err)
		}

		dr, err := algorithms.NewDeterministicReversible(dek)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to create deterministic reversible algo instance",
				logger.Err(err),
				logger.Base64("dekB64", dek))
			return nil, fmt.Errorf("failed to create deterministic reversible algo instance")
		}

		res := dr.Tokenize(ctx, pars.Plaintext)
		if res.Ciphertext == nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to tokenize with deterministic reversible algo",
				slog.String("plaintext", string(pars.Plaintext)),
				logger.Base64("dekB64", dek))
			return nil, fmt.Errorf("failed to tokenize with deterministic reversible algo")
		}

		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"successfully tokenized",
			logger.Base64("plaintext", pars.Plaintext),
			logger.Base64("wrappedDek", wrappedDek),
			logger.Base64("ciphertext", res.Ciphertext),
			slog.Bool("reversible", reversible),
			slog.Bool("deterministic", deterministic),
		)
		res.DekWrapped = wrappedDek

		return res, nil
	}

	return nil, fmt.Errorf("wrong tokenize params")
}

func (t *TokenizerService) Detokenize(ctx context.Context, pars *domain.DetokenizeParams) ([]byte, error) {
	deterministic := pars.Deterministic
	ciphertext := pars.Ciphertext
	wrappedDek := pars.WrappedDek

	var dek []byte
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(dek)

	if deterministic {
		dek, err := t.vault.UnwrapDEK(ctx, wrappedDek, t.convergentKey)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to unwrap DEK",
				logger.Base64("wrapped_dek", wrappedDek),
				slog.String("key", t.convergentKey))
			return nil, fmt.Errorf("failed to unwrap DEK: %w", err)
		}
		dr, err := algorithms.NewDeterministicReversible(dek)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to create deterministic reversible algo instance",
				slog.String("dek", string(dek)),
				logger.Err(err))
			return nil, fmt.Errorf("failed to create deterministic reversible algo instance")
		}

		res, err := dr.Detokenize(ctx, ciphertext)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to detokenize with deterministic reversible algo",
				logger.Err(err),
				logger.Base64("ciphertext", ciphertext),
				logger.Base64("dek", dek))
			return nil, fmt.Errorf("failed to detokenize with deterministic reversible algo")
		}

		return res, nil
	}

	return nil, fmt.Errorf("wrong detokenize params")
}
