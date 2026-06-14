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

// isExpired reports whether mapping has outlived its TTL. A TokenTtl of 0 means the
// mapping never expires.
func isExpired(mapping *domain.Mapping) bool {
	return mapping.TokenTtl != 0 && mapping.CreatedAt.Add(mapping.TokenTtl).Before(time.Now())
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

	if mapping != nil && isExpired(mapping) {
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

		if isExpired(mapping) {
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

func (m *MappingService) GetMappingByToken(ctx context.Context, token string) (*domain.Mapping, error) {
	mapping, err := m.storage.SelectMappingByToken(ctx, token)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to get mapping by token from storage",
			logger.Err(err))
		return nil, err
	}

	if isExpired(mapping) {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "mapping expired",
			slog.String("id", mapping.ID.String()))
		if err = m.storage.DeleteMappingById(ctx, mapping.ID); err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from storage",
				slog.String("id", mapping.ID.String()),
				logger.Err(err))
		}
		if err = m.cache.DeleteMappingById(ctx, mapping.ID); err != nil {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to delete mapping by id from cache",
				slog.String("id", mapping.ID.String()),
				logger.Err(err))
		}
		return nil, errs.ErrMappingExpired
	}

	if err = m.cache.SaveMapping(ctx, mapping, m.cacheTtl); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "failed to save mapping in cache",
			slog.String("id", mapping.ID.String()),
			logger.Err(err))
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

func (m *MappingService) CreateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error) {
	result, err := m.storage.CreateKind(ctx, kind)
	if err != nil {
		if errors.Is(err, errs.ErrKindAlreadyExists) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "kind already exists")
			return nil, err
		}

		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to insert kind",
			logger.Err(err))

		return nil, err
	}

	return result, nil
}

func (m *MappingService) GetKindById(ctx context.Context, id int32) (*domain.Kind, error) {
	kind, err := m.storage.GetKindById(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to get kind by id",
			slog.Int("id", int(id)),
			logger.Err(err))
		return nil, err
	}

	return kind, nil
}

func (m *MappingService) GetKindByName(ctx context.Context, name string) (*domain.Kind, error) {
	kind, err := m.storage.GetKindByName(ctx, name)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to get kind by name",
			slog.String("name", name),
			logger.Err(err))
		return nil, err
	}

	return kind, nil
}

func (m *MappingService) GetAllKinds(ctx context.Context) ([]*domain.Kind, error) {
	kinds, err := m.storage.GetAllKinds(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to get kinds list",
			logger.Err(err))
		return nil, err
	}

	return kinds, nil
}

func (m *MappingService) UpdateKind(ctx context.Context, kind *domain.Kind) (*domain.Kind, error) {
	updatedKind, err := m.storage.UpdateKind(ctx, kind)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to update kind",
			slog.Int("id", int(kind.Id)),
			logger.Err(err))
		return nil, err
	}

	return updatedKind, nil
}

func (m *MappingService) DeleteKindById(ctx context.Context, id int32) error {
	err := m.storage.DeleteKindById(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to delete kind",
			slog.Int("id", int(id)),
			logger.Err(err))
		return err
	}

	return nil
}

func (m *MappingService) CreateAuditLog(ctx context.Context, entry *domain.AuditLogEntry) (*domain.AuditLogEntry, error) {
	result, err := m.storage.CreateAuditLog(ctx, entry)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to insert audit log entry",
			logger.Err(err))
		return nil, err
	}

	return result, nil
}

func (m *MappingService) GetAuditLogList(ctx context.Context) ([]*domain.AuditLogEntry, error) {
	entries, err := m.storage.GetAuditLogList(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to get audit log list",
			logger.Err(err))
		return nil, err
	}

	return entries, nil
}

func (m *MappingService) UpdateMappingDek(ctx context.Context, id uuid.UUID, dekWrapped []byte) error {
	if err := m.storage.UpdateMappingDek(ctx, id, dekWrapped); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to update mapping dek",
			slog.String("id", id.String()),
			logger.Err(err))
		return err
	}

	if err := m.cache.DeleteMappingById(ctx, id); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to delete mapping from cache",
			slog.String("id", id.String()),
			logger.Err(err))
	}

	return nil
}

func (m *MappingService) UpdateMappingCrypto(ctx context.Context, id uuid.UUID, dekWrapped, cipherText []byte, algoName string) error {
	if err := m.storage.UpdateMappingCrypto(ctx, id, dekWrapped, cipherText, algoName); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to update mapping crypto",
			slog.String("id", id.String()),
			logger.Err(err))
		return err
	}

	if err := m.cache.DeleteMappingById(ctx, id); err != nil {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"failed to delete mapping from cache",
			slog.String("id", id.String()),
			logger.Err(err))
	}

	return nil
}
