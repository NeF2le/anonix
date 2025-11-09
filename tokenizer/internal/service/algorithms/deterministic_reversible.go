package algorithms

import (
	"context"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/miscreant/miscreant.go"
	"log/slog"
)

type DeterministicReversible struct {
	aead cipher.AEAD
}

func NewDeterministicReversible(dek []byte) (*DeterministicReversible, error) {
	var aead cipher.AEAD
	var err error
	ctx := context.Background()
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Creating AEAD with DEK", slog.String("dek", base64.StdEncoding.EncodeToString(dek)))

	aead, err = miscreant.NewAEAD("AES-SIV", dek, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	return &DeterministicReversible{aead: aead}, nil
}

func (s *DeterministicReversible) Tokenize(ctx context.Context, plaintext []byte) *domain.TokenResult {
	ciphertext := s.aead.Seal(nil, nil, plaintext, nil)

	res := &domain.TokenResult{
		Ciphertext: ciphertext,
		AlgoName:   "aes-256-siv",
	}

	return res
}

func (s *DeterministicReversible) Detokenize(ctx context.Context, ciphertext []byte) ([]byte, error) {
	plaintext, err := s.aead.Open(nil, nil, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open ciphertext: %w", err)
	}

	return plaintext, nil
}
