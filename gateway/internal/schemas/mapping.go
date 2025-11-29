package schemas

import (
	"time"
)

type ExistsTokenSchema struct {
	Token string `json:"token"`
}

type CreateMappingSchema struct {
	Token         string        `json:"token"`
	Ciphertext    []byte        `json:"ciphertext"`
	DekWrapped    []byte        `json:"dek_wrapped"`
	TypeId        int32         `json:"type_id"`
	ExpiresAt     time.Time     `json:"expires_at"`
	Deterministic bool          `json:"deterministic"`
	Reversible    bool          `json:"reversible"`
	CacheTTL      time.Duration `json:"cache_ttl"`
}

type UpdateMappingSchema struct {
	TokenTtl time.Duration `json:"token_ttl"`
}

type MappingSchema struct {
	Id            string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CipherText    string `json:"cipher_text,omitempty" example:"rfAsPo1N0a3cBiELOlgmHUXS5Q=="`
	DekWrapped    string `json:"dek_wrapped,omitempty" example:"dmF1bHQ6djE6dTZCY0lXRURFOG..."`
	TokenTtl      string `json:"token_ttl,omitempty" example:"24h0m0s"`
	CreatedAt     string `json:"created_at,omitempty" example:"2006-01-02T15:04:05Z07:00"`
	Deterministic bool   `json:"deterministic,omitempty" example:"true"`
	Reversible    bool   `json:"reversible,omitempty" example:"true"`
}
