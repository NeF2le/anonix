package main

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/grpc/runner"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/NeF2le/anonix/mapping/internal/config"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/cache"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/storage"
	"github.com/NeF2le/anonix/mapping/internal/service"
	transportgrpc "github.com/NeF2le/anonix/mapping/internal/transport/grpc"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
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

	redisClient, err := redis.NewRedisClient(ctx, &cfg.Redis, cfg.Mapping.RedisDB)
	if err != nil {
		panic(err)
	}
	redisServerCfg := redis.GetRedisConfigInfo(ctx, redisClient)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "redis server config", slog.String("config", redisServerCfg))

	postgresClient, err := postgres.NewPostgresClient(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}
	err = postgres.Migrate(ctx, &cfg.Postgres, cfg.MigrationsPath)
	if err != nil {
		log.Fatal(err)
	}

	cacheAdapter := cache.NewRedisAdapter(redisClient)
	storageAdapter := storage.NewPostgresAdapter(postgresClient)

	mappingService := service.NewMappingService(
		storageAdapter,
		cacheAdapter,
		cfg.Mapping.CacheTtl,
	)
	grpcHandler := transportgrpc.NewGRPCMappingHandler(mappingService)
	grpcServer, err := transportgrpc.CreateGRPC(ctx, grpcHandler)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating grpc mapping server: %w", err))
	}

	go runner.MustRunGRPC(ctx, grpcServer, cfg.Mapping.Port, cfg.Mapping.Host)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "mapping shutting down")
}
