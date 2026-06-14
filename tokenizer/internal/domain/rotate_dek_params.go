package domain

type RotateDEKParams struct {
	WrappedDek    []byte
	Ciphertext    []byte
	Deterministic bool
	AlgoName      string
}
