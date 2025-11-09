package ports

import (
	"context"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"time"
)

type StorageRepository interface {
	SelectMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	SelectAllMappings(ctx context.Context) ([]*domain.Mapping, error)
	InsertMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error)
	UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error)
	DeleteMappingById(ctx context.Context, id uuid.UUID) error
}

type CacheRepository interface {
	GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	SaveMapping(ctx context.Context, mapping *domain.Mapping, ttl time.Duration) error
	DeleteMappingById(ctx context.Context, id uuid.UUID) error
}

type MappingUseCase interface {
	GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error)
	GetAllMappings(ctx context.Context) ([]*domain.Mapping, error)
	CreateMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error)
	UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error)
	DeleteMappingById(ctx context.Context, id uuid.UUID) error
}
