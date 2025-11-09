package service

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/ports"
	"log/slog"
)

type MappingCleanerService struct {
	storage ports.StorageRepository
	cache   ports.CacheRepository
}

func NewMappingCleanerService(storage ports.StorageRepository, cache ports.CacheRepository) *MappingCleanerService {
	return &MappingCleanerService{
		storage: storage,
		cache:   cache,
	}
}

func (m *MappingCleanerService) DeleteExpiredMappings(ctx context.Context) (int, int, error) {
	var storageCount, cacheCount int

	ids, err := m.storage.DeleteExpiredMappings(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Error deleting expired mappings",
			logger.Err(err))
		return 0, 0, fmt.Errorf("delete expired mappings error: %w", err)
	}

	storageCount = len(ids)

	for _, id := range ids {
		err = m.cache.DeleteMappingById(ctx, id)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "Error deleting mapping from cache",
				slog.String("id", id.String()),
				logger.Err(err))
		}
		cacheCount++
	}

	return storageCount, cacheCount, nil
}
