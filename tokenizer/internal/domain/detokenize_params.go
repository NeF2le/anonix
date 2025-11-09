package domain

type DetokenizeParams struct {
	Ciphertext    []byte
	WrappedDek    []byte
	Deterministic bool
}
