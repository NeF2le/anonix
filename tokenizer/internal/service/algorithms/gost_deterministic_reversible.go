package algorithms

import (
	"context"
	"crypto/cipher"
	"crypto/hmac"
	"fmt"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/pedroalbanese/gogost/gost34112012256"
	"github.com/pedroalbanese/gogost/gost3412128"
	"github.com/pedroalbanese/gogost/mgm"
)

type GostDeterministicReversible struct {
	aead cipher.AEAD
	dek  []byte
}

func NewGostDeterministicReversible(dek []byte) (*GostDeterministicReversible, error) {
	blockCipher := gost3412128.NewCipher(dek)

	aead, err := mgm.NewMGM(blockCipher, gostNonceSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	return &GostDeterministicReversible{aead: aead, dek: dek}, nil
}

func (s *GostDeterministicReversible) Tokenize(ctx context.Context, plaintext []byte) *domain.TokenResult {
	nonce := s.deriveNonce(plaintext)

	ciphertext := s.aead.Seal(nil, nonce, plaintext, nil)

	return &domain.TokenResult{
		Ciphertext: append(nonce, ciphertext...),
		AlgoName:   "gost-kuznechik-mgm",
	}
}

func (s *GostDeterministicReversible) Detokenize(ctx context.Context, ciphertext []byte) ([]byte, error) {
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

// deriveNonce computes a synthetic, deterministic MGM nonce from the DEK and plaintext,
// analogous to AES-SIV's synthetic IV: identical plaintexts produce identical nonces
// (and therefore identical ciphertexts), while different plaintexts yield different nonces.
func (s *GostDeterministicReversible) deriveNonce(plaintext []byte) []byte {
	mac := hmac.New(gost34112012256.New, s.dek)
	mac.Write(plaintext)

	nonce := mac.Sum(nil)[:gostNonceSize]
	nonce[0] &= 0x7F

	return nonce
}
