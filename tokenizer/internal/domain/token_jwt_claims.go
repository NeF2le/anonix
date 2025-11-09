package domain

type TokenClaims struct {
	WrappedDek    []byte
	Ciphertext    []byte
	Deterministic bool
	Reversible    bool
}
