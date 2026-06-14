package domain

import (
	"github.com/google/uuid"
	"time"
)

type Mapping struct {
	ID            uuid.UUID     `json:"id"`
	Token         string        `json:"token"`
	DekWrapped    []byte        `json:"dek_wrapped,omitempty"`
	CipherText    []byte        `json:"cipher_text,omitempty"`
	TokenTtl      time.Duration `json:"token_ttl,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	Deterministic bool          `json:"deterministic"`
	Kind          *Kind         `json:"kind"`
	AlgoName      string        `json:"algo_name"`
}
