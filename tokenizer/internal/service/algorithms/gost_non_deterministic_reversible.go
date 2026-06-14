package algorithms

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/pedroalbanese/gogost/gost3412128"
	"github.com/pedroalbanese/gogost/mgm"
)

// gostNonceSize is the size, in bytes, of the MGM nonce/tag for the Kuznechik (128-bit block) cipher.
const gostNonceSize = gost3412128.BlockSize

type GostNonDeterministicReversible struct {
	aead cipher.AEAD
}

func NewGostNonDeterministicReversible(dek []byte) (*GostNonDeterministicReversible, error) {
	blockCipher := gost3412128.NewCipher(dek)

	aead, err := mgm.NewMGM(blockCipher, gostNonceSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	return &GostNonDeterministicReversible{aead: aead}, nil
}

func (s *GostNonDeterministicReversible) Tokenize(ctx context.Context, plaintext []byte) *domain.TokenResult {
	nonce := make([]byte, gostNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return &domain.TokenResult{}
	}
	nonce[0] &= 0x7F

	ciphertext := s.aead.Seal(nil, nonce, plaintext, nil)

	return &domain.TokenResult{
		Ciphertext: append(nonce, ciphertext...),
		AlgoName:   "gost-kuznechik-mgm-random",
	}
}

func (s *GostNonDeterministicReversible) Detokenize(ctx context.Context, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < gostNonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ct := ciphertext[:gostNonceSize], ciphertext[gostNonceSize:]

	plaintext, err := s.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open ciphertext: %w", err)
	}

	return plaintext, nil
}
