package algorithms

import "crypto/rand"

// NonDeterministicSuffix generates a short, random token suffix unrelated to the
// input plaintext: repeated calls with the same plaintext yield different suffixes.
type NonDeterministicSuffix struct{}

func NewNonDeterministicSuffix() *NonDeterministicSuffix {
	return &NonDeterministicSuffix{}
}

func (n *NonDeterministicSuffix) GenerateSuffix(plaintext []byte) []byte {
	out := make([]byte, tokenSuffixSize)
	if _, err := rand.Read(out); err != nil {
		panic(err)
	}

	return out
}
