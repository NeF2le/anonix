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

func (m *grpcMappingHandler) GetMappingByToken(ctx context.Context, req *mapping.GetMappingByTokenRequest) (
	*mapping.GetMappingResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	mappingOut, err := m.mapping.GetMappingByToken(ctx, req.GetToken())
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

func (m *grpcMappingHandler) CreateKind(ctx context.Context, req *mapping.CreateKindRequest) (
	*mapping.CreateKindResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	kindIn := helpers.CreateKindRequestToModel(req)

	kindOut, err := m.mapping.CreateKind(ctx, kindIn)
	if err != nil {
		if errors.Is(err, errs.ErrKindAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "kind already exists")
		}
		return nil, status.Error(codes.Internal, "failed to insert kind")
	}

	return &mapping.CreateKindResponse{
		Kind: helpers.ModelToGRPCKind(kindOut),
	}, nil
}

func (m *grpcMappingHandler) GetKind(ctx context.Context, req *mapping.GetKindRequest) (
	*mapping.GetKindResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	kindOut, err := m.mapping.GetKindById(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, errs.ErrKindNotFound) {
			return nil, status.Error(codes.NotFound, "kind not found")
		}
		return nil, status.Error(codes.Internal, "failed to get kind")
	}

	return &mapping.GetKindResponse{
		Kind: helpers.ModelToGRPCKind(kindOut),
	}, nil
}

func (m *grpcMappingHandler) GetKindByName(ctx context.Context, req *mapping.GetKindByNameRequest) (
	*mapping.GetKindByNameResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	kindOut, err := m.mapping.GetKindByName(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, errs.ErrKindNotFound) {
			return nil, status.Error(codes.NotFound, "kind not found")
		}
		return nil, status.Error(codes.Internal, "failed to get kind")
	}

	return &mapping.GetKindByNameResponse{
		Kind: helpers.ModelToGRPCKind(kindOut),
	}, nil
}

func (m *grpcMappingHandler) ListKinds(ctx context.Context, req *mapping.ListKindsRequest) (
	*mapping.ListKindsResponse, error) {
	kinds, err := m.mapping.GetAllKinds(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get kinds")
	}

	var kindModels []*mapping.Kind
	for _, kind := range kinds {
		kindModels = append(kindModels, helpers.ModelToGRPCKind(kind))
	}

	return &mapping.ListKindsResponse{
		Kinds: kindModels,
	}, nil
}

func (m *grpcMappingHandler) UpdateKind(ctx context.Context, req *mapping.UpdateKindRequest) (
	*mapping.UpdateKindResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	kindIn := helpers.UpdateKindRequestToModel(req)

	kindOut, err := m.mapping.UpdateKind(ctx, kindIn)
	if err != nil {
		if errors.Is(err, errs.ErrKindNotFound) {
			return nil, status.Error(codes.NotFound, "kind not found")
		}
		return nil, status.Error(codes.Internal, "failed to update kind")
	}

	return &mapping.UpdateKindResponse{
		Kind: helpers.ModelToGRPCKind(kindOut),
	}, nil
}

func (m *grpcMappingHandler) DeleteKind(ctx context.Context, req *mapping.DeleteKindRequest) (
	*mapping.DeleteKindResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := m.mapping.DeleteKindById(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, errs.ErrKindNotFound) {
			return nil, status.Error(codes.NotFound, "kind not found")
		}
		if errors.Is(err, errs.ErrKindInUse) {
			return nil, status.Error(codes.FailedPrecondition, "kind is in use")
		}
		return nil, status.Error(codes.Internal, "failed to delete kind")
	}

	return &mapping.DeleteKindResponse{}, nil
}

func (m *grpcMappingHandler) CreateAuditLog(ctx context.Context, req *mapping.CreateAuditLogRequest) (
	*mapping.CreateAuditLogResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetAction() == "" {
		return nil, status.Error(codes.InvalidArgument, "action is required")
	}

	entry, err := helpers.CreateAuditLogRequestToModel(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	result, err := m.mapping.CreateAuditLog(ctx, entry)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create audit log entry")
	}

	return &mapping.CreateAuditLogResponse{Entry: helpers.ModelToGRPCAuditLogEntry(result)}, nil
}

func (m *grpcMappingHandler) GetAuditLogList(ctx context.Context, req *mapping.GetAuditLogListRequest) (
	*mapping.GetAuditLogListResponse, error) {
	entries, err := m.mapping.GetAuditLogList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get audit log list")
	}

	var result []*mapping.AuditLogEntry
	for _, entry := range entries {
		result = append(result, helpers.ModelToGRPCAuditLogEntry(entry))
	}

	return &mapping.GetAuditLogListResponse{Entries: result}, nil
}

func (m *grpcMappingHandler) UpdateMappingDek(ctx context.Context, req *mapping.UpdateMappingDekRequest) (
	*mapping.UpdateMappingDekResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetDekWrapped() == nil {
		return nil, status.Error(codes.InvalidArgument, "dek wrapping is required")
	}

	mappingUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "mapping id is invalid")
	}

	if err = m.mapping.UpdateMappingDek(ctx, mappingUUID, req.GetDekWrapped()); err != nil {
		if errors.Is(err, errs.ErrMappingNotFound) {
			return nil, status.Error(codes.NotFound, "mapping not found")
		}
		return nil, status.Error(codes.Internal, "failed to update mapping dek")
	}

	return &mapping.UpdateMappingDekResponse{}, nil
}

func (m *grpcMappingHandler) UpdateMappingCrypto(ctx context.Context, req *mapping.UpdateMappingCryptoRequest) (
	*mapping.UpdateMappingCryptoResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetDekWrapped() == nil {
		return nil, status.Error(codes.InvalidArgument, "dek wrapping is required")
	}
	if req.GetCipherText() == nil {
		return nil, status.Error(codes.InvalidArgument, "cipher text is required")
	}

	mappingUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "mapping id is invalid")
	}

	if err = m.mapping.UpdateMappingCrypto(ctx, mappingUUID, req.GetDekWrapped(), req.GetCipherText(), req.GetAlgoName()); err != nil {
		if errors.Is(err, errs.ErrMappingNotFound) {
			return nil, status.Error(codes.NotFound, "mapping not found")
		}
		return nil, status.Error(codes.Internal, "failed to update mapping crypto")
	}

	return &mapping.UpdateMappingCryptoResponse{}, nil
}
