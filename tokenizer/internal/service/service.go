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
	deterministic, pseudonymize := pars.Deterministic, pars.Pseudonymize

	var suffixAlgo algorithms.TokenSuffixAlgorithm
	if deterministic {
		suffixAlgo = algorithms.NewDeterministicSuffix([]byte(t.convergentKey))
	} else {
		suffixAlgo = algorithms.NewNonDeterministicSuffix()
	}

	res := &domain.TokenResult{
		TokenSuffix: suffixAlgo.GenerateSuffix(pars.Plaintext),
	}

	if pseudonymize {
		wrappedDek, dek, err := t.vault.GenerateDEK(ctx, t.dekBitsLength, t.convergentKey)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				"failed to generate DEK",
				slog.Int("dek_bits_length", t.dekBitsLength),
				slog.String("key", t.convergentKey),
				logger.Err(err),
			)
			return nil, fmt.Errorf("failed to generate DEK: %w", err)
		}
		defer func(b []byte) {
			for i := range b {
				b[i] = 0
			}
		}(dek)

		algo, err := newAlgoForTokenize(pars.Algorithm, deterministic, dek)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx,
				"failed to create algo instance",
				slog.Bool("deterministic", deterministic),
				slog.String("algorithm", pars.Algorithm),
				logger.Err(err))
			return nil, fmt.Errorf("failed to create algo instance: %w", err)
		}

		cipherRes := algo.Tokenize(ctx, pars.Plaintext)
		res.Ciphertext = cipherRes.Ciphertext
		res.AlgoName = cipherRes.AlgoName
		res.DekWrapped = wrappedDek
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx,
		"successfully tokenized",
		logger.Base64("tokenSuffix", res.TokenSuffix),
		slog.Bool("pseudonymize", pseudonymize),
		slog.Bool("deterministic", deterministic),
		slog.String("algo", res.AlgoName),
	)

	return res, nil
}

func (t *TokenizerService) Detokenize(ctx context.Context, pars *domain.DetokenizeParams) ([]byte, error) {
	dek, err := t.vault.UnwrapDEK(ctx, pars.WrappedDek, t.convergentKey)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to unwrap DEK",
			slog.String("key", t.convergentKey),
			logger.Err(err))
		return nil, fmt.Errorf("failed to unwrap DEK: %w", err)
	}
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(dek)

	algo, err := newAlgoForDetokenize(pars.AlgoName, pars.Deterministic, dek)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to create algo instance",
			slog.Bool("deterministic", pars.Deterministic),
			slog.String("algo_name", pars.AlgoName),
			logger.Err(err))
		return nil, fmt.Errorf("failed to create algo instance: %w", err)
	}

	res, err := algo.Detokenize(ctx, pars.Ciphertext)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to detokenize",
			logger.Err(err))
		return nil, fmt.Errorf("failed to detokenize")
	}

	return res, nil
}

func (t *TokenizerService) RotateMasterKey(ctx context.Context) error {
	if err := t.vault.RotateKey(ctx, t.convergentKey); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to rotate master key",
			slog.String("key", t.convergentKey),
			logger.Err(err))
		return fmt.Errorf("failed to rotate master key: %w", err)
	}
	return nil
}

func (t *TokenizerService) RewrapDEK(ctx context.Context, wrappedDek []byte) ([]byte, error) {
	newWrappedDek, err := t.vault.RewrapDEK(ctx, wrappedDek, t.convergentKey)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to rewrap DEK",
			slog.String("key", t.convergentKey),
			logger.Err(err))
		return nil, fmt.Errorf("failed to rewrap DEK: %w", err)
	}
	return newWrappedDek, nil
}

func (t *TokenizerService) RotateDEK(ctx context.Context, pars *domain.RotateDEKParams) (*domain.RotateDEKResult, error) {
	oldDek, err := t.vault.UnwrapDEK(ctx, pars.WrappedDek, t.convergentKey)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to unwrap DEK",
			slog.String("key", t.convergentKey),
			logger.Err(err))
		return nil, fmt.Errorf("failed to unwrap DEK: %w", err)
	}
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(oldDek)

	decryptAlgo, err := newAlgoForDetokenize(pars.AlgoName, pars.Deterministic, oldDek)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to create algo instance",
			slog.Bool("deterministic", pars.Deterministic),
			slog.String("algo_name", pars.AlgoName),
			logger.Err(err))
		return nil, fmt.Errorf("failed to create algo instance: %w", err)
	}

	plaintext, err := decryptAlgo.Detokenize(ctx, pars.Ciphertext)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to detokenize",
			logger.Err(err))
		return nil, fmt.Errorf("failed to detokenize")
	}
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(plaintext)

	newWrappedDek, newDek, err := t.vault.GenerateDEK(ctx, t.dekBitsLength, t.convergentKey)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to generate DEK",
			slog.Int("dek_bits_length", t.dekBitsLength),
			slog.String("key", t.convergentKey),
			logger.Err(err))
		return nil, fmt.Errorf("failed to generate DEK: %w", err)
	}
	defer func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}(newDek)

	encryptAlgo, err := newAlgoForTokenize(algoFamily(pars.AlgoName), pars.Deterministic, newDek)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"failed to create algo instance",
			slog.Bool("deterministic", pars.Deterministic),
			slog.String("algo_name", pars.AlgoName),
			logger.Err(err))
		return nil, fmt.Errorf("failed to create algo instance: %w", err)
	}

	cipherRes := encryptAlgo.Tokenize(ctx, plaintext)

	return &domain.RotateDEKResult{
		DekWrapped: newWrappedDek,
		Ciphertext: cipherRes.Ciphertext,
		AlgoName:   cipherRes.AlgoName,
	}, nil
}

// algoFamily maps a persisted AlgoName back to the "algorithm" selector
// accepted by newAlgoForTokenize, so DEK rotation re-encrypts with the same
// algorithm family the data was originally encrypted with.
func algoFamily(algoName string) string {
	switch algoName {
	case "gost-kuznechik-mgm", "gost-kuznechik-mgm-random":
		return "gost-kuznechik"
	default:
		return "aes-siv"
	}
}

// newAlgoForTokenize selects the algorithm implementation used to encrypt
// new data, based on the requested algorithm family and determinism mode.
func newAlgoForTokenize(algorithm string, deterministic bool, dek []byte) (algorithms.Algorithm, error) {
	switch algorithm {
	case "", "aes-siv":
		if deterministic {
			return algorithms.NewDeterministicReversible(dek)
		}
		return algorithms.NewNonDeterministicReversible(dek)
	case "gost-kuznechik":
		if deterministic {
			return algorithms.NewGostDeterministicReversible(dek)
		}
		return algorithms.NewGostNonDeterministicReversible(dek)
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}
}

// newAlgoForDetokenize selects the algorithm implementation used to decrypt
// existing data, based on the persisted AlgoName.
func newAlgoForDetokenize(algoName string, deterministic bool, dek []byte) (algorithms.Algorithm, error) {
	switch algoName {
	case "":
		// Legacy mappings created before algorithm selection was introduced - always AES-SIV.
		if deterministic {
			return algorithms.NewDeterministicReversible(dek)
		}
		return algorithms.NewNonDeterministicReversible(dek)
	case "aes-256-siv":
		return algorithms.NewDeterministicReversible(dek)
	case "aes-256-siv-random":
		return algorithms.NewNonDeterministicReversible(dek)
	case "gost-kuznechik-mgm":
		return algorithms.NewGostDeterministicReversible(dek)
	case "gost-kuznechik-mgm-random":
		return algorithms.NewGostNonDeterministicReversible(dek)
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algoName)
	}
}
