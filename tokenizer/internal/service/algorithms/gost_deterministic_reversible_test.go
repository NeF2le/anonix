package algorithms

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"
)

func TestGostDeterministicReversible_TokenizeDetokenize(t *testing.T) {
	ctx := context.Background()
	dek := make([]byte, dekLength)
	if _, err := rand.Read(dek); err != nil {
		t.Fatalf("failed to generate dek: %v", err)
	}
	s, err := NewGostDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new gost deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res.Ciphertext == nil {
		t.Fatalf("Tokenize returned nil ciphertext")
	}
	if res.AlgoName != "gost-kuznechik-mgm" {
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

func TestGostDeterministicReversible_Deterministic(t *testing.T) {
	ctx := context.Background()
	dek := make([]byte, dekLength)
	if _, err := rand.Read(dek); err != nil {
		t.Fatalf("failed to generate dek: %v", err)
	}
	s, err := NewGostDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new gost deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res1 := s.Tokenize(ctx, plaintext)
	res2 := s.Tokenize(ctx, plaintext)

	if !bytes.Equal(res1.Ciphertext, res2.Ciphertext) {
		t.Fatalf("Tokenize returned different ciphertexts for the same input: %x vs %x", res1.Ciphertext, res2.Ciphertext)
	}
}
