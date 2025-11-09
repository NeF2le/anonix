package domain

type TokenizeParams struct {
	Plaintext     []byte
	Deterministic bool
	Reversible    bool
}
