package domain

type TokenResult struct {
	TokenSuffix []byte
	Ciphertext  []byte
	DekWrapped  []byte
	AlgoName    string
	KeyName     string
}
