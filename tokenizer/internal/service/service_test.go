package service

import (
	"bytes"
	"context"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/miscreant/miscreant.go"
	"testing"
)

const testConvergentKey = "test-convergent-key"
const testDekBitsLength = 256

type fakeVault struct{}

func (f *fakeVault) GenerateDEK(ctx context.Context, bits int, keyName string) ([]byte, []byte, error) {
	dek := miscreant.GenerateKey(bits / 8)
	wrapped := make([]byte, len(dek))
	copy(wrapped, dek)
	return wrapped, dek, nil
}

func (f *fakeVault) UnwrapDEK(ctx context.Context, wrappedDek []byte, keyName string) ([]byte, error) {
	return wrappedDek, nil
}

func (f *fakeVault) RotateKey(ctx context.Context, keyName string) error {
	return nil
}

func (f *fakeVault) RewrapDEK(ctx context.Context, wrappedDek []byte, keyName string) ([]byte, error) {
	return wrappedDek, nil
}

func TestTokenizerService_Tokenize_AllCombinations(t *testing.T) {
	svc := NewTokenizerService(&fakeVault{}, testConvergentKey, testDekBitsLength)
	ctx := context.Background()
	plaintext := []byte("Корнилов Евгений Александрович")

	combos := []struct {
		name          string
		deterministic bool
		pseudonymize  bool
	}{
		{"anonymize_deterministic", true, false},
		{"anonymize_nondeterministic", false, false},
		{"pseudonymize_deterministic", true, true},
		{"pseudonymize_nondeterministic", false, true},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			res, err := svc.Tokenize(ctx, &domain.TokenizeParams{
				Plaintext:     plaintext,
				Deterministic: c.deterministic,
				Pseudonymize:  c.pseudonymize,
			})
			if err != nil {
				t.Fatalf("Tokenize returned error: %v", err)
			}
			if len(res.TokenSuffix) == 0 {
				t.Fatalf("Tokenize returned empty token suffix")
			}

			if c.pseudonymize {
				if res.Ciphertext == nil {
					t.Fatalf("expected non-nil Ciphertext for pseudonymized token")
				}
				if res.DekWrapped == nil {
					t.Fatalf("expected non-nil DekWrapped for pseudonymized token")
				}

				plainOut, err := svc.Detokenize(ctx, &domain.DetokenizeParams{
					Ciphertext:    res.Ciphertext,
					WrappedDek:    res.DekWrapped,
					Deterministic: c.deterministic,
				})
				if err != nil {
					t.Fatalf("Detokenize returned error: %v", err)
				}
				if !bytes.Equal(plainOut, plaintext) {
					t.Fatalf("decrypted plaintext mismatch: got=%q want=%q", string(plainOut), string(plaintext))
				}
			} else {
				if res.Ciphertext != nil {
					t.Fatalf("expected nil Ciphertext for anonymized token, got %x", res.Ciphertext)
				}
				if res.DekWrapped != nil {
					t.Fatalf("expected nil DekWrapped for anonymized token, got %x", res.DekWrapped)
				}
			}

			res2, err := svc.Tokenize(ctx, &domain.TokenizeParams{
				Plaintext:     plaintext,
				Deterministic: c.deterministic,
				Pseudonymize:  c.pseudonymize,
			})
			if err != nil {
				t.Fatalf("second Tokenize returned error: %v", err)
			}

			if c.deterministic {
				if !bytes.Equal(res.TokenSuffix, res2.TokenSuffix) {
					t.Fatalf("expected deterministic token suffix to match: got=%x want=%x", res2.TokenSuffix, res.TokenSuffix)
				}
			} else {
				if bytes.Equal(res.TokenSuffix, res2.TokenSuffix) {
					t.Fatalf("expected non-deterministic token suffixes to differ, both=%x", res.TokenSuffix)
				}
			}
		})
	}
}
