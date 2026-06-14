package domain

type RotateDEKResult struct {
	DekWrapped []byte
	Ciphertext []byte
	AlgoName   string
}
