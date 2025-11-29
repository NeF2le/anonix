package helpers

import (
	"encoding/base64"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"time"
)

func ProtoMappingToSchema(m *mapping.MappingModel) *schemas.MappingSchema {
	ttl := ""
	if m.TokenTtl != nil {
		ttl = m.TokenTtl.AsDuration().String()
	}

	return &schemas.MappingSchema{
		Id:            m.Id,
		CipherText:    base64.StdEncoding.EncodeToString(m.CipherText),
		DekWrapped:    base64.StdEncoding.EncodeToString(m.DekWrapped),
		Deterministic: m.Deterministic,
		Reversible:    m.Reversible,
		TokenTtl:      ttl,
		CreatedAt:     m.CreatedAt.AsTime().Format(time.RFC3339),
	}
}
