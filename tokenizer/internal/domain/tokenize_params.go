package domain

type TokenizeParams struct {
	Plaintext     []byte
	Deterministic bool
	Pseudonymize  bool
	Algorithm     string
}
