package algorithms

import (
	"crypto/hmac"
	"crypto/sha256"
)

// DeterministicSuffix generates a short token suffix from plaintext: the same
// plaintext always yields the same suffix, while different plaintexts produce
// different suffixes with overwhelming probability.
type DeterministicSuffix struct {
	key []byte
}

func NewDeterministicSuffix(key []byte) *DeterministicSuffix {
	return &DeterministicSuffix{key: key}
}

func (d *DeterministicSuffix) GenerateSuffix(plaintext []byte) []byte {
	mac := hmac.New(sha256.New, d.key)
	mac.Write(plaintext)

	return mac.Sum(nil)[:tokenSuffixSize]
}
