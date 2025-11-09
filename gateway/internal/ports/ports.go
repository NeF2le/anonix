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
}

type MappingServiceRepository interface {
	GetMapping(ctx context.Context, req *mapping.GetMappingRequest) (*mapping.GetMappingResponse, error)
	GetMappingList(ctx context.Context, req *mapping.GetMappingListRequest) (*mapping.GetMappingListResponse, error)
	CreateMapping(ctx context.Context, req *mapping.CreateMappingRequest) (*mapping.CreateMappingResponse, error)
	DeleteMapping(ctx context.Context, req *mapping.DeleteMappingRequest) (*mapping.DeleteMappingResponse, error)
	UpdateMapping(ctx context.Context, req *mapping.UpdateMappingRequest) (*mapping.UpdateMappingResponse, error)
}

type AuthServiceRepository interface {
	Register(ctx context.Context, req *auth_service.RegisterRequest) (*auth_service.RegisterResponse, error)
	Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error)
	Refresh(ctx context.Context, req *auth_service.RefreshRequest) (*auth_service.RefreshResponse, error)
	IsAdmin(ctx context.Context, req *auth_service.IsAdminRequest) (*auth_service.IsAdminResponse, error)
}
