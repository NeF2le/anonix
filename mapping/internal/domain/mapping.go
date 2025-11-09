package domain

import (
	"github.com/google/uuid"
	"time"
)

type Mapping struct {
	ID            uuid.UUID     `json:"id"`
	DekWrapped    []byte        `json:"dek_wrapped,omitempty"`
	CipherText    []byte        `json:"cipher_text,omitempty"`
	TokenTtl      time.Duration `json:"token_ttl,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	Deterministic bool          `json:"deterministic"`
	Reversible    bool          `json:"reversible"`
}
