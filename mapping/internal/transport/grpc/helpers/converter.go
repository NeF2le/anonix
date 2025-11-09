package helpers

import (
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func CreateMappingRequestToModel(req *mapping.CreateMappingRequest) *domain.Mapping {
	m := &domain.Mapping{
		CipherText:    req.CipherText,
		DekWrapped:    req.DekWrapped,
		Reversible:    req.Reversible,
		Deterministic: req.Deterministic,
	}

	var tokenTtl time.Duration
	if ttl := req.GetTokenTtl(); ttl != nil {
		tokenTtl = req.GetTokenTtl().AsDuration()
	}
	m.TokenTtl = tokenTtl

	return m
}

func GPRCMappingToModel(model *mapping.MappingModel) *domain.Mapping {
	mappingUUID, _ := uuid.Parse(model.Id)

	return &domain.Mapping{
		ID:            mappingUUID,
		DekWrapped:    model.DekWrapped,
		Reversible:    model.Reversible,
		Deterministic: model.Deterministic,
		CipherText:    model.CipherText,
		TokenTtl:      model.TokenTtl.AsDuration(),
		CreatedAt:     model.CreatedAt.AsTime(),
	}
}

func ModelToGRPCMapping(model *domain.Mapping) *mapping.MappingModel {
	return &mapping.MappingModel{
		Id:            model.ID.String(),
		DekWrapped:    model.DekWrapped,
		Reversible:    model.Reversible,
		Deterministic: model.Deterministic,
		CipherText:    model.CipherText,
		TokenTtl:      durationpb.New(model.TokenTtl),
		CreatedAt:     timestamppb.New(model.CreatedAt),
	}
}
