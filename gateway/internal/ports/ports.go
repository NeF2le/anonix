package ports

import (
	"context"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
)

type TokenizerServiceRepository interface {
	Tokenize(ctx context.Context, req *tokenizer.TokenizeRequest) (*tokenizer.TokenizeResponse, error)
	Detokenize(ctx context.Context, req *tokenizer.DetokenizeRequest) (*tokenizer.DetokenizeResponse, error)
	RotateMasterKey(ctx context.Context, req *tokenizer.RotateMasterKeyRequest) (*tokenizer.RotateMasterKeyResponse, error)
	RewrapDEK(ctx context.Context, req *tokenizer.RewrapDEKRequest) (*tokenizer.RewrapDEKResponse, error)
	RotateDEK(ctx context.Context, req *tokenizer.RotateDEKRequest) (*tokenizer.RotateDEKResponse, error)
}

type MappingServiceRepository interface {
	GetMapping(ctx context.Context, req *mapping.GetMappingRequest) (*mapping.GetMappingResponse, error)
	GetMappingByToken(ctx context.Context, req *mapping.GetMappingByTokenRequest) (*mapping.GetMappingResponse, error)
	GetMappingList(ctx context.Context, req *mapping.GetMappingListRequest) (*mapping.GetMappingListResponse, error)
	CreateMapping(ctx context.Context, req *mapping.CreateMappingRequest) (*mapping.CreateMappingResponse, error)
	DeleteMapping(ctx context.Context, req *mapping.DeleteMappingRequest) (*mapping.DeleteMappingResponse, error)
	UpdateMapping(ctx context.Context, req *mapping.UpdateMappingRequest) (*mapping.UpdateMappingResponse, error)
	UpdateMappingDek(ctx context.Context, req *mapping.UpdateMappingDekRequest) (*mapping.UpdateMappingDekResponse, error)
	UpdateMappingCrypto(ctx context.Context, req *mapping.UpdateMappingCryptoRequest) (*mapping.UpdateMappingCryptoResponse, error)

	GetKind(ctx context.Context, req *mapping.GetKindRequest) (*mapping.GetKindResponse, error)
	GetKindByName(ctx context.Context, req *mapping.GetKindByNameRequest) (*mapping.GetKindByNameResponse, error)
	ListKinds(ctx context.Context, req *mapping.ListKindsRequest) (*mapping.ListKindsResponse, error)
	CreateKind(ctx context.Context, req *mapping.CreateKindRequest) (*mapping.CreateKindResponse, error)
	UpdateKind(ctx context.Context, req *mapping.UpdateKindRequest) (*mapping.UpdateKindResponse, error)
	DeleteKind(ctx context.Context, req *mapping.DeleteKindRequest) (*mapping.DeleteKindResponse, error)

	CreateAuditLog(ctx context.Context, req *mapping.CreateAuditLogRequest) (*mapping.CreateAuditLogResponse, error)
	GetAuditLogList(ctx context.Context, req *mapping.GetAuditLogListRequest) (*mapping.GetAuditLogListResponse, error)
}

type AuthServiceRepository interface {
	Register(ctx context.Context, req *auth_service.RegisterRequest) (*auth_service.RegisterResponse, error)
	Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error)
	Refresh(ctx context.Context, req *auth_service.RefreshRequest) (*auth_service.RefreshResponse, error)
	IsAdmin(ctx context.Context, req *auth_service.IsAdminRequest) (*auth_service.IsAdminResponse, error)

	GetUsers(ctx context.Context, req *auth_service.GetUsersRequest) (*auth_service.GetUsersResponse, error)
	DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*auth_service.DeleteUserResponse, error)

	AssignRole(ctx context.Context, req *auth_service.AssignRoleRequest) (*auth_service.AssignRoleResponse, error)
	RemoveRole(ctx context.Context, req *auth_service.RemoveRoleRequest) (*auth_service.RemoveRoleResponse, error)
	UpdateClearanceLevel(ctx context.Context, req *auth_service.UpdateClearanceLevelRequest) (*auth_service.UpdateClearanceLevelResponse, error)
	GetRolesList(ctx context.Context, req *auth_service.GetRolesListRequest) (*auth_service.GetRolesListResponse, error)
	GetUserRoles(ctx context.Context, req *auth_service.GetUserRolesRequest) (*auth_service.GetUserRolesResponse, error)
}
