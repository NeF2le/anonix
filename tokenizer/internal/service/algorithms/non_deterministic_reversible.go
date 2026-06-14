package algorithms

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/miscreant/miscreant.go"
)

// sivADSize is the size, in bytes, of the random associated data prepended to the
// ciphertext to make AES-SIV non-deterministic while remaining reversible.
const sivADSize = 16

type NonDeterministicReversible struct {
	aead cipher.AEAD
}

func NewNonDeterministicReversible(dek []byte) (*NonDeterministicReversible, error) {
	aead, err := miscreant.NewAEAD("AES-SIV", dek, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	return &NonDeterministicReversible{aead: aead}, nil
}

func (s *NonDeterministicReversible) Tokenize(ctx context.Context, plaintext []byte) *domain.TokenResult {
	ad := make([]byte, sivADSize)
	if _, err := rand.Read(ad); err != nil {
		return &domain.TokenResult{}
	}

	ciphertext := s.aead.Seal(nil, nil, plaintext, ad)

	return &domain.TokenResult{
		Ciphertext: append(ad, ciphertext...),
		AlgoName:   "aes-256-siv-random",
	}
}

func (s *NonDeterministicReversible) Detokenize(ctx context.Context, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < sivADSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	ad, ct := ciphertext[:sivADSize], ciphertext[sivADSize:]

	plaintext, err := s.aead.Open(nil, nil, ct, ad)
	if err != nil {
		return nil, fmt.Errorf("failed to open ciphertext: %w", err)
	}

	return plaintext, nil
}
