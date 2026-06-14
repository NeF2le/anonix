package ports

import (
	"context"
	"github.com/NeF2le/anonix/mapping/internal/domain"
)

type VaultRepository interface {
	GenerateDEK(ctx context.Context, bits int, keyName string) ([]byte, []byte, error)
	UnwrapDEK(ctx context.Context, wrappedDek []byte, keyName string) ([]byte, error)
	RotateKey(ctx context.Context, keyName string) error
	RewrapDEK(ctx context.Context, wrappedDek []byte, keyName string) ([]byte, error)
}

type TokenizerUseCase interface {
	Tokenize(ctx context.Context, pars *domain.TokenizeParams) (*domain.TokenResult, error)
	Detokenize(ctx context.Context, pars *domain.DetokenizeParams) ([]byte, error)
	RotateMasterKey(ctx context.Context) error
	RewrapDEK(ctx context.Context, wrappedDek []byte) ([]byte, error)
	RotateDEK(ctx context.Context, pars *domain.RotateDEKParams) (*domain.RotateDEKResult, error)
}
