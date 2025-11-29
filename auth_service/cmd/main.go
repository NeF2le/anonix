package main

import (
	"context"
	"github.com/NeF2le/anonix/auth_service/internal/config"
	"github.com/NeF2le/anonix/auth_service/internal/ports/adapters/cache"
	"github.com/NeF2le/anonix/auth_service/internal/ports/adapters/storage"
	"github.com/NeF2le/anonix/auth_service/internal/service"
	transportgrpc "github.com/NeF2le/anonix/auth_service/internal/transport/grpc"
	"github.com/NeF2le/anonix/common/grpc/runner"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), logger.KeyForLogLevel, cfg.LogLevel)
	ctx = logger.New(ctx)

	postgresClient, err := postgres.NewPostgresClient(ctx, &cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := redis.NewRedisClient(ctx, &cfg.Redis, cfg.AuthService.RedisDB)
	if err != nil {
		log.Fatal(err)
	}

	err = postgres.Migrate(ctx, &cfg.Postgres, cfg.MigrationsPath)
	if err != nil {
		log.Fatal(err)
	}

	storageAdapter := storage.NewAuthPostgresAdapter(postgresClient)
	cacheAdapter := cache.NewAuthRedisAdapter(redisClient, cfg.RefreshExpiration, cfg.AccessExpiration)

	authService := service.NewAuthService(
		storageAdapter,
		cacheAdapter,
		cfg.JwtSecret,
		cfg.RefreshExpiration,
		cfg.AccessExpiration,
	)

	grpcHandler := transportgrpc.NewGRPCAuthHandler(authService)

	var grpcServer *grpc.Server
	tlsCfg := cfg.TLS
	if tlsCfg.Enabled {
		if err = tls_helpers.Verification(cfg.AuthService.Host, &tlsCfg); err != nil {
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
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	go runner.MustRunGRPC(ctx, grpcServer, cfg.AuthService.Port, cfg.AuthService.Host)

	<-ctx.Done()
	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "auth service stopped")
}
