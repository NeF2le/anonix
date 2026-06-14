package ports

import (
	"context"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"time"
)

type StorageRepository interface {
	SelectMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	SelectMappingByToken(ctx context.Context, token string) (*domain.Mapping, error)
	SelectAllMappings(ctx context.Context) ([]*domain.Mapping, error)
	InsertMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error)
	UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error)
	UpdateMappingDek(ctx context.Context, id uuid.UUID, dekWrapped []byte) error
	UpdateMappingCrypto(ctx context.Context, id uuid.UUID, dekWrapped, cipherText []byte, algoName string) error
	DeleteMappingById(ctx context.Context, id uuid.UUID) error

	GetKindById(ctx context.Context, id int32) (*domain.Kind, error)
	GetKindByName(ctx context.Context, name string) (*domain.Kind, error)
	GetAllKinds(ctx context.Context) ([]*domain.Kind, error)
	CreateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error)
	UpdateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error)
	DeleteKindById(ctx context.Context, id int32) error

	CreateAuditLog(ctx context.Context, entry *domain.AuditLogEntry) (*domain.AuditLogEntry, error)
	GetAuditLogList(ctx context.Context) ([]*domain.AuditLogEntry, error)
}

type CacheRepository interface {
	GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	SaveMapping(ctx context.Context, mapping *domain.Mapping, ttl time.Duration) error
	DeleteMappingById(ctx context.Context, id uuid.UUID) error
}

type MappingUseCase interface {
	GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	GetMappingByToken(ctx context.Context, token string) (*domain.Mapping, error)
	GetAllMappings(ctx context.Context) ([]*domain.Mapping, error)
	CreateMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error)
	UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error)
	UpdateMappingDek(ctx context.Context, id uuid.UUID, dekWrapped []byte) error
	UpdateMappingCrypto(ctx context.Context, id uuid.UUID, dekWrapped, cipherText []byte, algoName string) error
	DeleteMappingById(ctx context.Context, id uuid.UUID) error

	GetKindById(ctx context.Context, id int32) (*domain.Kind, error)
	GetKindByName(ctx context.Context, name string) (*domain.Kind, error)
	GetAllKinds(ctx context.Context) ([]*domain.Kind, error)
	CreateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error)
	UpdateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error)
	DeleteKindById(ctx context.Context, id int32) error

	CreateAuditLog(ctx context.Context, entry *domain.AuditLogEntry) (*domain.AuditLogEntry, error)
	GetAuditLogList(ctx context.Context) ([]*domain.AuditLogEntry, error)
}
