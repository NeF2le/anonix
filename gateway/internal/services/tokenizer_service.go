package services

import (
	"context"
	"github.com/NeF2le/anonix/common/callers"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/gateway/internal/ports"
	"time"
)

type TokenizerService struct {
	TokenizerServiceRepo ports.TokenizerServiceRepository
	MaxRetries           uint
	BaseDelay            time.Duration
}

func NewTokenizerService(
	tokenizerServiceRepo ports.TokenizerServiceRepository,
	maxRetries uint,
	baseDelay time.Duration) *TokenizerService {
	return &TokenizerService{
		TokenizerServiceRepo: tokenizerServiceRepo,
		MaxRetries:           maxRetries,
		BaseDelay:            baseDelay,
	}
}

func (t *TokenizerService) Tokenize(ctx context.Context, req *tokenizer.TokenizeRequest) (
	*tokenizer.TokenizeResponse, error) {
	resultChan := make(chan *tokenizer.TokenizeResponse, 1)

	err := callers.Retry(func() error {
		resp, err := t.TokenizerServiceRepo.Tokenize(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, t.MaxRetries, t.BaseDelay)

	if err != nil {
		return nil, err
	}

	return <-resultChan, nil
}

func (t *TokenizerService) Detokenize(ctx context.Context, req *tokenizer.DetokenizeRequest) (
	*tokenizer.DetokenizeResponse, error) {
	resultChan := make(chan *tokenizer.DetokenizeResponse, 1)

	err := callers.Retry(func() error {
		resp, err := t.TokenizerServiceRepo.Detokenize(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, t.MaxRetries, t.BaseDelay)

	if err != nil {
		return nil, err
	}

	return <-resultChan, nil
}
