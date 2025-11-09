package main

import (
	"context"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/NeF2le/anonix/mapping/internal/config"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/cache"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/storage"
	"github.com/NeF2le/anonix/mapping/internal/service"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	ctx = logger.New(ctx)

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, logger.KeyForLogLevel, cfg.LogLevel)

	redisClient, err := redis.NewRedisClient(ctx, &cfg.Redis, cfg.MappingCleaner.RedisDB)
	if err != nil {
		panic(err)
	}

	postgresClient, err := postgres.NewPostgresClient(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	cacheAdapter := cache.NewRedisAdapter(redisClient)
	storageAdapter := storage.NewPostgresAdapter(postgresClient)

	mappingCleanerService := service.NewMappingCleanerService(storageAdapter, cacheAdapter)

	ticker := time.NewTicker(cfg.MappingCleaner.Cooldown)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.GetLoggerFromCtx(ctx).Info(ctx, "mapping shutting down")
			return
		case <-ticker.C:
			storageDel, cacheDel, delErr := mappingCleanerService.DeleteExpiredMappings(ctx)
			if delErr != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx, "error deleting expired mappings", logger.Err(delErr))
			}
			logger.GetLoggerFromCtx(ctx).Info(ctx, "deleted expired mapping ",
				slog.Int("from storage", storageDel),
				slog.Int("from cache", cacheDel))
		}
	}
}
