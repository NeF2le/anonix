package domain

type TokenResult struct {
	Ciphertext []byte
	DekWrapped []byte
	AlgoName   string
	KeyName    string
}
