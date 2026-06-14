package algorithms

import (
	"context"
	"github.com/NeF2le/anonix/mapping/internal/domain"
)

// tokenSuffixSize is the size, in bytes, of the short token suffix shown to the user.
const tokenSuffixSize = 4

type Algorithm interface {
	Tokenize(ctx context.Context, plaintext []byte) *domain.TokenResult
	Detokenize(ctx context.Context, ciphertext []byte) ([]byte, error)
}

// TokenSuffixAlgorithm generates the short suffix used to build a user-facing token,
// e.g. "fio_7f82a1c3".
type TokenSuffixAlgorithm interface {
	GenerateSuffix(plaintext []byte) []byte
}
