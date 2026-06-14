package helpers

import (
	"encoding/base64"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"time"
)

func ProtoMappingToSchema(m *mapping.MappingModel) *schemas.MappingSchema {
	ttl := ""
	if m.TokenTtl != nil {
		ttl = m.TokenTtl.AsDuration().String()
	}

	result := &schemas.MappingSchema{
		Id:            m.Id,
		Token:         m.Token,
		CipherText:    base64.StdEncoding.EncodeToString(m.CipherText),
		DekWrapped:    base64.StdEncoding.EncodeToString(m.DekWrapped),
		Deterministic: m.Deterministic,
		TokenTtl:      ttl,
		CreatedAt:     m.CreatedAt.AsTime().Format(time.RFC3339),
		AlgoName:      m.AlgoName,
	}

	if m.Kind != nil {
		result.Kind = ProtoKindToSchema(m.Kind)
	}

	return result
}

func ProtoRoleToSchema(r *auth_service.Role) *schemas.RoleSchema {
	return &schemas.RoleSchema{
		Id:   r.Id,
		Name: r.Name,
	}
}

func ProtoUserToSchema(u *auth_service.User) *schemas.UserSchema {
	roles := make([]*schemas.RoleSchema, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, ProtoRoleToSchema(r))
	}

	return &schemas.UserSchema{
		Id:             u.Id,
		Login:          u.Login,
		Roles:          roles,
		ClearanceLevel: u.ClearanceLevel,
	}
}

func ProtoKindToSchema(k *mapping.Kind) *schemas.KindSchema {
	if k == nil {
		return nil
	}

	return &schemas.KindSchema{
		Id:          k.Id,
		Name:        k.Name,
		RussianName: k.RussianName,
		AccessLevel: k.AccessLevel,
		Mask:        k.Mask,
		ShortName:   k.ShortName,
	}
}

func ProtoAuditLogEntryToSchema(e *mapping.AuditLogEntry) *schemas.AuditLogEntrySchema {
	result := &schemas.AuditLogEntrySchema{
		Id:        e.Id,
		UserId:    e.UserId,
		Action:    e.Action,
		Token:     e.Token,
		CreatedAt: e.CreatedAt.AsTime().Format(time.RFC3339),
	}

	if e.Kind != nil {
		result.Kind = ProtoKindToSchema(e.Kind)
	}

	return result
}
