package main

import (
	"context"
	"github.com/NeF2le/anonix/common/grpc/runner"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/NeF2le/anonix/common/vault_agent"
	"github.com/NeF2le/anonix/mapping/internal/config"
	"github.com/NeF2le/anonix/mapping/internal/ports/adapters/vault"
	"github.com/NeF2le/anonix/mapping/internal/service"
	transportgrpc "github.com/NeF2le/anonix/mapping/internal/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(logger.New(ctx), logger.KeyForLogLevel, cfg.LogLevel)

	vaultAgent, err := vault_agent.NewVaultAgent(ctx, &cfg.VaultAgent)
	if err != nil {
		panic(err)
	}
	hashicorpAdapter := vault.NewHashiCorpAdapter(vaultAgent)

	tokenizerService := service.NewTokenizerService(hashicorpAdapter, cfg.ConvergentKey, cfg.DEKBitsLength)
	grpcHandler := transportgrpc.NewGRPCTokenizerHandler(tokenizerService)

	var grpcServer *grpc.Server
	tlsCfg := cfg.TLS
	if tlsCfg.Enabled {
		if err = tls_helpers.Verification(cfg.Tokenizer.Host, &tlsCfg); err != nil {
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
		panic(err)
	}

	go runner.MustRunGRPC(ctx, grpcServer, cfg.Tokenizer.Port, cfg.Tokenizer.Host)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "tokenizer shutting down")
}
