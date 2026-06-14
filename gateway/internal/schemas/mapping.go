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
	KindId        int           `json:"kind_id"`
	CacheTTL      time.Duration `json:"cache_ttl"`
}

type UpdateMappingSchema struct {
	TokenTtl time.Duration `json:"token_ttl"`
}

type MappingSchema struct {
	Id            string      `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Token         string      `json:"token" example:"fio_7f82a1c3"`
	CipherText    string      `json:"cipher_text,omitempty" example:"rfAsPo1N0a3cBiELOlgmHUXS5Q=="`
	DekWrapped    string      `json:"dek_wrapped,omitempty" example:"dmF1bHQ6djE6dTZCY0lXRURFOG..."`
	TokenTtl      string      `json:"token_ttl,omitempty" example:"24h0m0s"`
	CreatedAt     string      `json:"created_at,omitempty" example:"2006-01-02T15:04:05Z07:00"`
	Deterministic bool        `json:"deterministic,omitempty" example:"true"`
	Kind          *KindSchema `json:"kind,omitempty"`
	AlgoName      string      `json:"algo_name,omitempty" example:"aes-256-siv"`
}

type CreateKindSchema struct {
	Name        string `json:"name" example:"passport"`
	RussianName string `json:"russian_name" example:"Паспорт"`
	AccessLevel int32  `json:"access_level" example:"3"`
	Mask        string `json:"mask" example:"^\\d{4} \\d{6}$"`
	ShortName   string `json:"short_name" example:"psp"`
}

type UpdateKindSchema struct {
	Name        string `json:"name" example:"passport"`
	RussianName string `json:"russian_name" example:"Паспорт"`
	AccessLevel int32  `json:"access_level" example:"3"`
	Mask        string `json:"mask" example:"^\\d{4} \\d{6}$"`
	ShortName   string `json:"short_name" example:"psp"`
}

type KindSchema struct {
	Id          int32  `json:"id" example:"1"`
	Name        string `json:"name" example:"passport"`
	RussianName string `json:"russian_name" example:"Паспорт"`
	AccessLevel int32  `json:"access_level" example:"3"`
	Mask        string `json:"mask" example:"^\\d{4} \\d{6}$"`
	ShortName   string `json:"short_name" example:"psp"`
}

type AuditLogEntrySchema struct {
	Id        string      `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserId    string      `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Action    string      `json:"action" example:"tokenize"`
	Token     string      `json:"token" example:"fio_7f82a1c3"`
	Kind      *KindSchema `json:"kind,omitempty"`
	CreatedAt string      `json:"created_at" example:"2006-01-02T15:04:05Z07:00"`
}
