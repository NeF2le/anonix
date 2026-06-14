package algorithms

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"
)

func TestGostNonDeterministicReversible_TokenizeDetokenize(t *testing.T) {
	ctx := context.Background()
	dek := make([]byte, dekLength)
	if _, err := rand.Read(dek); err != nil {
		t.Fatalf("failed to generate dek: %v", err)
	}
	s, err := NewGostNonDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new gost non-deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res.Ciphertext == nil {
		t.Fatalf("Tokenize returned nil ciphertext")
	}
	if res.AlgoName != "gost-kuznechik-mgm-random" {
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

func TestGostNonDeterministicReversible_NonDeterministic(t *testing.T) {
	ctx := context.Background()
	dek := make([]byte, dekLength)
	if _, err := rand.Read(dek); err != nil {
		t.Fatalf("failed to generate dek: %v", err)
	}
	s, err := NewGostNonDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new gost non-deterministic reversible: %v", err)
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
