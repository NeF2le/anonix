package algorithms

import (
	"bytes"
	"context"
	"github.com/miscreant/miscreant.go"
	"testing"
)

func TestNonDeterministicReversible_TokenizeDetokenize(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewNonDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new non-deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res.Ciphertext == nil {
		t.Fatalf("Tokenize returned nil ciphertext")
	}
	if res.AlgoName != "aes-256-siv-random" {
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

func TestNonDeterministicReversible_NonDeterministic(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewNonDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new non-deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res1 := s.Tokenize(ctx, plaintext)
	res2 := s.Tokenize(ctx, plaintext)

	if bytes.Equal(res1.Ciphertext, res2.Ciphertext) {
		t.Fatalf("Tokenize returned identical ciphertexts for the same input: %x", res1.Ciphertext)
	}

	for _, res := range []*struct{ ct []byte }{{res1.Ciphertext}, {res2.Ciphertext}} {
		plainOut, err := s.Detokenize(ctx, res.ct)
		if err != nil {
			t.Fatalf("Detokenize returned error: %v", err)
		}
		if !bytes.Equal(plainOut, plaintext) {
			t.Fatalf("decrypted plaintext mismatch: got=%q want=%q", string(plainOut), string(plaintext))
		}
	}
}
