package algorithms

import (
	"bytes"
	"context"
	"github.com/miscreant/miscreant.go"
	"testing"
)

// dekLength is shared across this package's tests: 32 bytes is valid both for
// AES-CMAC-SIV (32 or 64) and for FF3 (AES-128/192/256 keys).
const dekLength = 32

func TestDeterministicReversible_TokenizeDetokenize(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res.Ciphertext == nil {
		t.Fatalf("Tokenize returned nil ciphertext")
	}
	if res.AlgoName != "aes-256-siv" {
		t.Fatalf("unexpected AlgoName: %s", res.AlgoName)
	}

	plainOut, err := s.Detokenize(ctx, res.Ciphertext)
	if err != nil {
		t.Fatalf("Detokenize returned error: %v", err)
	}
	if !bytes.Equal(plainOut, plaintext) {
		t.Fatalf("decrypted plaintext mismatch: got=%q want=%q", string(plainOut), string(plaintext))
	}
}

func TestDeterministicReversible_Deterministic(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res1 := s.Tokenize(ctx, plaintext)
	res2 := s.Tokenize(ctx, plaintext)

	if !bytes.Equal(res1.Ciphertext, res2.Ciphertext) {
		t.Fatalf("Tokenize returned different ciphertexts for the same input: %x vs %x", res1.Ciphertext, res2.Ciphertext)
	}
}
