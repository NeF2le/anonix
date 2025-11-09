package ports

import (
	"context"
	"github.com/google/uuid"
)

type StorageRepository interface {
	DeleteExpiredMappings(ctx context.Context) ([]uuid.UUID, error)
}

type CacheRepository interface {
	DeleteMappingById(ctx context.Context, id uuid.UUID) error
}
