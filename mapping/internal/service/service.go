package service

import (
	"context"
	"errors"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/NeF2le/anonix/mapping/internal/ports"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type MappingService struct {
	storage  ports.StorageRepository
	cache    ports.CacheRepository
	cacheTtl time.Duration
}

func NewMappingService(
	storage ports.StorageRepository,
	cache ports.CacheRepository,
	cacheTTL time.Duration,
) *MappingService {
	return &MappingService{
		storage:  storage,
		cache:    cache,
		cacheTtl: cacheTTL,
	}
}

func (m *MappingService) CreateMapping(ctx context.Context, mapping *domain.Mapping) (*domain.Mapping, error) {
	result, err := m.storage.InsertMapping(ctx, mapping)
	if err != nil {
		if errors.Is(err, errs.ErrMappingAlreadyExists) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "mapping already exists")
			return nil, err
		}
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to insert mapping in storage", logger.Err(err))
		return nil, err
	}

	if err = m.cache.SaveMapping(ctx, mapping, m.cacheTtl); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to save mapping in cache",
			slog.String("id", mapping.ID.String()),
			logger.Err(err))
	} else {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "successfully saved mapping in cache",
			slog.String("id", mapping.ID.String()),
			slog.String("mapping created at", mapping.CreatedAt.Format("2006-01-02 15:04:05")),
			slog.String("mapping ttl", mapping.TokenTtl.String()))
	}

	return result, nil
}

func (m *MappingService) DeleteMappingById(ctx context.Context, id uuid.UUID) error {
	err := m.storage.DeleteMappingById(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping from storage",
			slog.String("id", id.String()),
			logger.Err(err))
		return err
	}
	err = m.cache.DeleteMappingById(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping from cache",
			slog.String("id", id.String()),
			logger.Err(err))
		return err
	}
	return nil
}

func (m *MappingService) GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error) {
	mapping, err := m.cache.GetMappingById(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to get mapping by id from cache", logger.Err(err))
	}

	if mapping != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "check mapping TTL",
			slog.String("id", id.String()),
			slog.String("now", time.Now().Format("2006-01-02 15:04:05")),
			slog.String("created at", mapping.CreatedAt.Format("2006-01-02 15:04:05")),
			slog.String("ttl", mapping.TokenTtl.String()),
			slog.String("time to delete", mapping.CreatedAt.Add(mapping.TokenTtl).Format("2006-01-02 15:04:05")))
	}

	if mapping != nil && mapping.CreatedAt.Add(mapping.TokenTtl).Before(time.Now()) {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "mapping expired",
			slog.String("id", id.String()))
		err = m.storage.DeleteMappingById(ctx, id)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from storage",
				slog.String("id", id.String()),
				logger.Err(err))
		}
		err = m.cache.DeleteMappingById(ctx, id)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from cache",
				slog.String("id", id.String()),
				logger.Err(err))
		}
		return nil, errs.ErrMappingExpired
	}

	if mapping == nil {
		mapping, err = m.storage.SelectMappingById(ctx, id)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to get mapping by id from storage",
				slog.String("id", id.String()),
				logger.Err(err))
			return nil, err
		}

		if mapping.CreatedAt.Add(mapping.TokenTtl).Before(time.Now()) {
			logger.GetLoggerFromCtx(ctx).Debug(ctx, "mapping expired",
				slog.String("id", id.String()))
			err = m.storage.DeleteMappingById(ctx, id)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Debug(ctx, "failed to delete mapping by id from storage",
					slog.String("id", id.String()),
					logger.Err(err))
			}
			err = m.cache.DeleteMappingById(ctx, id)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from cache",
					slog.String("id", id.String()),
					logger.Err(err))
			}
			return nil, errs.ErrMappingExpired
		}

		if err = m.cache.DeleteMappingById(ctx, id); err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from cache",
				slog.String("id", id.String()),
				logger.Err(err))
		}
		if err = m.cache.SaveMapping(ctx, mapping, m.cacheTtl); err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to save mapping in cache",
				slog.String("id", id.String()),
				logger.Err(err))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "saved mapping in cache",
				slog.String("id", id.String()))
		}
	}

	return mapping, nil
}

func (m *MappingService) GetAllMappings(ctx context.Context) ([]*domain.Mapping, error) {
	mappings, err := m.storage.SelectAllMappings(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to get mapping list", logger.Err(err))
		return nil, err
	}
	return mappings, nil
}

func (m *MappingService) UpdateMapping(ctx context.Context, id uuid.UUID, tokenTtl time.Duration) (*domain.Mapping, error) {
	mapping, err := m.storage.UpdateMapping(ctx, id, tokenTtl)
	if err != nil {
		return nil, err
	}
	if err = m.cache.SaveMapping(ctx, mapping, m.cacheTtl); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to save mapping in cache",
			slog.String("id", id.String()),
			logger.Err(err))
	} else {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "successfully saved mapping in cache",
			slog.String("id", id.String()))
	}

	return mapping, nil
}
