package main

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/grpc/runner"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/NeF2le/anonix/mapping/internal/config"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/cache"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/storage"
	"github.com/NeF2le/anonix/mapping/internal/service"
	transportgrpc "github.com/NeF2le/anonix/mapping/internal/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	var grpcServer *grpc.Server
	tlsCfg := cfg.TLS
	if tlsCfg.Enabled {
		if err = tls_helpers.Verification(cfg.Mapping.Host, &tlsCfg); err != nil {
			panic(err)
		}

		var grpcTls credentials.TransportCredentials
		grpcTls, err = tls_helpers.LoadServerTLSConfig(tlsCfg.ServerPublicKey, tlsCfg.ServerPrivateKey, tlsCfg.RootPublicKey)
		if err != nil {
			panic(err)
		}
		grpcServer, err = transportgrpc.CreateGRPCTLS(grpcHandler, grpcTls)
		logger.GetLoggerFromCtx(ctx).Info(ctx, "tokenizer grpc server created with tls")
	} else {
		grpcServer, err = transportgrpc.CreateGRPC(grpcHandler)
		logger.GetLoggerFromCtx(ctx).Info(ctx, "tokenizer grpc server created without tls")
	}

	if err != nil {
		log.Fatal(fmt.Errorf("error creating grpc mapping server: %w", err))
	}

	go runner.MustRunGRPC(ctx, grpcServer, cfg.Mapping.Port, cfg.Mapping.Host)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "mapping shutting down")
}
