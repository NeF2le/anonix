package schemas

type TokenizeSchema struct {
	Plaintext     []byte `json:"plaintext"`
	Deterministic bool   `json:"deterministic"`
	Mode          string `json:"mode" example:"pseudonymize"` // "pseudonymize" | "anonymize"
	TokenTTL      int64  `json:"token_ttl"`
	KindId        int    `json:"kind_id"`
	Algorithm     string `json:"algorithm" example:"aes-siv"` // "" | "aes-siv" | "gost-kuznechik"
}

type DetokenizeSchema struct {
	Token string `json:"token"`
}

type DetokenizeRespSchema struct {
	Plaintext []byte `json:"plaintext"`
}

// TokenizeResultSchema is returned for anonymized tokens, which are not stored by the
// mapping service and therefore have no id, ttl, etc.
type TokenizeResultSchema struct {
	Token string `json:"token" example:"fio_7f82a1c3"`
}
