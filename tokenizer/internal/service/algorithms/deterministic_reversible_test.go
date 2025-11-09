package algorithms

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/miscreant/miscreant.go"
	"testing"
)

const dekLength = 64

func TestTokenize_Success(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res == nil {
		t.Fatalf("Tokenize returned nil result")
	}
	if res.AlgoName != "aes-256-siv" {
		t.Fatalf("unexpected AlgoName: %s", res.AlgoName)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(res.Token)
	if err != nil {
		t.Fatalf("token is not valid base64.RawURLEncoding: %v", err)
	}
	if !bytes.Equal(decoded, res.Ciphertext) {
		t.Fatalf("token base64 does not match Ciphertext field: decoded len=%d ct len=%d", len(decoded), len(res.Ciphertext))
	}

	aead, err := miscreant.NewAEAD("AES-SIV", dek, 16)
	if err != nil {
		t.Fatalf("miscreant.NewAEAD: %v", err)
	}
	nonce := make([]byte, aead.NonceSize())

	plainOut, err := aead.Open(nil, nonce, res.Ciphertext, nil)
	if err != nil {
		t.Fatalf("AEAD.Open failed: %v", err)
	}
	if !bytes.Equal(plainOut, plaintext) {
		t.Fatalf("decrypted plaintext mismatch: got=%q want=%q", string(plainOut), string(plaintext))
	}
}

func TestTokenize_Deterministic(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res1 := s.Tokenize(ctx, plaintext)
	if res1.Token == "" {
		t.Fatalf("Tokenize returned empty token")
	}
	res2 := s.Tokenize(ctx, plaintext)
	if res2.Token == "" {
		t.Fatalf("Tokenize returned empty token")
	}

	if res1.Token != res2.Token {
		t.Fatalf("Tokenize returned different tokens: %q vs %q", res1.Token, res2.Token)
	}
	if !bytes.Equal(res1.Ciphertext, res2.Ciphertext) {
		t.Fatalf("Tokenize returned different ciphertexts: %q vs %q", res1.Ciphertext, res2.Ciphertext)
	}
}

func TestDetokenize_Success(t *testing.T) {
	ctx := context.Background()
	dek := miscreant.GenerateKey(dekLength)
	s, err := NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}

	plaintext := []byte("very strong secret string")
	res := s.Tokenize(ctx, plaintext)
	if res == nil {
		t.Fatalf("Tokenize returned nil result")
	}

	s, err = NewDeterministicReversible(dek)
	if err != nil {
		t.Fatalf("failed to create new deterministic reversible: %v", err)
	}
	resultPlaintext, err := s.Detokenize(ctx, res.Token)
	if err != nil {
		t.Fatalf("Detokenize returned error: %v", err)
	}

	if !bytes.Equal(resultPlaintext, plaintext) {
		t.Fatalf("Detokenize returned wrong result: got=%q want=%q", string(resultPlaintext), string(plaintext))
	}
}
