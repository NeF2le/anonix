package algorithms

import (
	"bytes"
	"testing"
)

func TestDeterministicSuffix_Deterministic(t *testing.T) {
	d := NewDeterministicSuffix([]byte("convergent-key"))

	plaintext := []byte("very strong secret string")
	suffix1 := d.GenerateSuffix(plaintext)
	suffix2 := d.GenerateSuffix(plaintext)

	if len(suffix1) != tokenSuffixSize {
		t.Fatalf("unexpected suffix length: got=%d want=%d", len(suffix1), tokenSuffixSize)
	}
	if !bytes.Equal(suffix1, suffix2) {
		t.Fatalf("GenerateSuffix returned different suffixes for the same input: %x vs %x", suffix1, suffix2)
	}
}

func TestDeterministicSuffix_DifferentInputsDifferentOutputs(t *testing.T) {
	d := NewDeterministicSuffix([]byte("convergent-key"))

	suffix1 := d.GenerateSuffix([]byte("plaintext one"))
	suffix2 := d.GenerateSuffix([]byte("plaintext two"))

	if bytes.Equal(suffix1, suffix2) {
		t.Fatalf("different plaintexts produced the same suffix: %x", suffix1)
	}
}
