package algorithms

import (
	"bytes"
	"testing"
)

func TestNonDeterministicSuffix_NonDeterministic(t *testing.T) {
	n := NewNonDeterministicSuffix()

	plaintext := []byte("very strong secret string")
	suffix1 := n.GenerateSuffix(plaintext)
	suffix2 := n.GenerateSuffix(plaintext)

	if len(suffix1) != tokenSuffixSize {
		t.Fatalf("unexpected suffix length: got=%d want=%d", len(suffix1), tokenSuffixSize)
	}
	if bytes.Equal(suffix1, suffix2) {
		t.Fatalf("GenerateSuffix returned identical suffixes for the same input: %x", suffix1)
	}
}
