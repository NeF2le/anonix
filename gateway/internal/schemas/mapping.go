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
