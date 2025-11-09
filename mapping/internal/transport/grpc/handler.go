package grpc

import (
	"context"
	"errors"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/mapping/internal/ports"
	"github.com/NeF2le/anonix/mapping/internal/transport/grpc/helpers"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcMappingHandler struct {
	mapping ports.MappingUseCase
	mapping.UnimplementedMappingServer
}

func NewGRPCMappingHandler(mapping ports.MappingUseCase) mapping.MappingServer {
	return &grpcMappingHandler{mapping: mapping}
}

func (m *grpcMappingHandler) CreateMapping(ctx context.Context, req *mapping.CreateMappingRequest) (*mapping.CreateMappingResponse, error) {
	if len(req.CipherText) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cipher text is required")
	}

	if req.Reversible && len(req.DekWrapped) == 0 {
		return nil, status.Error(codes.InvalidArgument, "token is reversible, wrapped dek is required")
	}

	mappingIn := helpers.CreateMappingRequestToModel(req)

	mappingOut, err := m.mapping.CreateMapping(ctx, mappingIn)
	if err != nil {
		if errors.Is(err, errs.ErrMappingAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "mapping already exists")
		}
		return nil, status.Error(codes.Internal, "failed to insert mapping")
	}

	return &mapping.CreateMappingResponse{MappingModel: helpers.ModelToGRPCMapping(mappingOut)}, nil
}

func (m *grpcMappingHandler) GetMapping(ctx context.Context, req *mapping.GetMappingRequest) (
	*mapping.GetMappingResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	mappingUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "mapping id is invalid")
	}

	mappingOut, err := m.mapping.GetMappingById(ctx, mappingUUID)
	if err != nil {
		if errors.Is(err, errs.ErrMappingNotFound) {
			return nil, status.Error(codes.NotFound, "mapping not found")
		}
		if errors.Is(err, errs.ErrMappingExpired) {
			return nil, status.Error(codes.DeadlineExceeded, "mapping expired")
		}
		return nil, status.Error(codes.Internal, "failed to get mapping")
	}

	return &mapping.GetMappingResponse{MappingModel: helpers.ModelToGRPCMapping(mappingOut)}, nil
}

func (m *grpcMappingHandler) DeleteMapping(ctx context.Context, req *mapping.DeleteMappingRequest) (
	*mapping.DeleteMappingResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	mappingUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "mapping id is invalid")
	}

	err = m.mapping.DeleteMappingById(ctx, mappingUUID)
	if err != nil {
		if errors.Is(err, errs.ErrMappingNotFound) {
			return nil, status.Error(codes.NotFound, "mapping not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete mapping")
	}

	return &mapping.DeleteMappingResponse{}, nil
}

func (m *grpcMappingHandler) GetMappingList(ctx context.Context, req *mapping.GetMappingListRequest) (
	*mapping.GetMappingListResponse, error) {
	mappings, err := m.mapping.GetAllMappings(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get mappings")
	}

	var mappingModels []*mapping.MappingModel
	for _, mp := range mappings {
		mappingModels = append(mappingModels, helpers.ModelToGRPCMapping(mp))
	}

	return &mapping.GetMappingListResponse{MappingModels: mappingModels}, err
}

func (m *grpcMappingHandler) UpdateMapping(ctx context.Context, req *mapping.UpdateMappingRequest) (
	*mapping.UpdateMappingResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	mappingUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "mapping id is invalid")
	}

	mappingOut, err := m.mapping.UpdateMapping(ctx, mappingUUID, req.TokenTtl.AsDuration())
	if err != nil {
		if errors.Is(err, errs.ErrMappingNotFound) {
			return nil, status.Error(codes.NotFound, "mapping not found")
		}
		return nil, status.Error(codes.Internal, "failed to update mapping")
	}

	return &mapping.UpdateMappingResponse{MappingModel: helpers.ModelToGRPCMapping(mappingOut)}, nil
}
