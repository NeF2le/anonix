package schemas

type TokenizeSchema struct {
	Plaintext     []byte `json:"plaintext"`
	Deterministic bool   `json:"deterministic"`
	Reversible    bool   `json:"reversible"`
	TokenTTL      int64  `json:"token_ttl"`
}

type DetokenizeSchema struct {
	Token string `json:"token"`
}

type DetokenizeRespSchema struct {
	Plaintext []byte `json:"plaintext"`
}
