package helpers

import (
	"fmt"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func CreateMappingRequestToModel(req *mapping.CreateMappingRequest) *domain.Mapping {
	m := &domain.Mapping{
		Token:         req.Token,
		CipherText:    req.CipherText,
		DekWrapped:    req.DekWrapped,
		Deterministic: req.Deterministic,
		AlgoName:      req.AlgoName,
	}

	var tokenTtl time.Duration
	if ttl := req.GetTokenTtl(); ttl != nil {
		tokenTtl = req.GetTokenTtl().AsDuration()
	}
	m.TokenTtl = tokenTtl

	if req.Kind != nil {
		m.Kind = GRPCKindToModel(req.Kind)
	}

	return m
}

func GPRCMappingToModel(model *mapping.MappingModel) *domain.Mapping {
	mappingUUID, _ := uuid.Parse(model.Id)

	m := &domain.Mapping{
		ID:            mappingUUID,
		Token:         model.Token,
		DekWrapped:    model.DekWrapped,
		Deterministic: model.Deterministic,
		CipherText:    model.CipherText,
		TokenTtl:      model.TokenTtl.AsDuration(),
		CreatedAt:     model.CreatedAt.AsTime(),
		AlgoName:      model.AlgoName,
	}

	if model.Kind != nil {
		m.Kind = GRPCKindToModel(model.Kind)
	}

	return m
}

func ModelToGRPCMapping(model *domain.Mapping) *mapping.MappingModel {
	m := &mapping.MappingModel{
		Id:            model.ID.String(),
		Token:         model.Token,
		DekWrapped:    model.DekWrapped,
		Deterministic: model.Deterministic,
		CipherText:    model.CipherText,
		TokenTtl:      durationpb.New(model.TokenTtl),
		CreatedAt:     timestamppb.New(model.CreatedAt),
		AlgoName:      model.AlgoName,
	}

	if model.Kind != nil {
		m.Kind = ModelToGRPCKind(model.Kind)
	}

	return m
}

func CreateKindRequestToModel(req *mapping.CreateKindRequest) *domain.Kind {
	return &domain.Kind{
		Name:        req.Name,
		RussianName: req.RussianName,
		AccessLevel: req.AccessLevel,
		Mask:        req.Mask,
		ShortName:   req.ShortName,
	}
}

func UpdateKindRequestToModel(req *mapping.UpdateKindRequest) *domain.Kind {
	return &domain.Kind{
		Id:          req.Id,
		Name:        req.Name,
		RussianName: req.RussianName,
		AccessLevel: req.AccessLevel,
		Mask:        req.Mask,
		ShortName:   req.ShortName,
	}
}

func GRPCKindToModel(kind *mapping.Kind) *domain.Kind {
	if kind == nil {
		return nil
	}

	return &domain.Kind{
		Id:          kind.Id,
		Name:        kind.Name,
		RussianName: kind.RussianName,
		AccessLevel: kind.AccessLevel,
		Mask:        kind.Mask,
		ShortName:   kind.ShortName,
	}
}

func ModelToGRPCKind(kind *domain.Kind) *mapping.Kind {
	if kind == nil {
		return nil
	}

	return &mapping.Kind{
		Id:          kind.Id,
		Name:        kind.Name,
		RussianName: kind.RussianName,
		AccessLevel: kind.AccessLevel,
		Mask:        kind.Mask,
		ShortName:   kind.ShortName,
	}
}

func CreateAuditLogRequestToModel(req *mapping.CreateAuditLogRequest) (*domain.AuditLogEntry, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	entry := &domain.AuditLogEntry{
		UserID: userID,
		Action: req.GetAction(),
		Token:  req.GetToken(),
	}

	if req.GetKindId() > 0 {
		entry.Kind = &domain.Kind{Id: req.GetKindId()}
	}

	return entry, nil
}

func ModelToGRPCAuditLogEntry(entry *domain.AuditLogEntry) *mapping.AuditLogEntry {
	e := &mapping.AuditLogEntry{
		Id:        entry.ID.String(),
		UserId:    entry.UserID.String(),
		Action:    entry.Action,
		Token:     entry.Token,
		CreatedAt: timestamppb.New(entry.CreatedAt),
	}

	if entry.Kind != nil {
		e.Kind = ModelToGRPCKind(entry.Kind)
	}

	return e
}
