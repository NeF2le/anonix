package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/NeF2le/anonix/gateway/internal/config"
	"github.com/NeF2le/anonix/gateway/internal/handlers/http_handlers"
	"github.com/NeF2le/anonix/gateway/internal/handlers/middlewares"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/auth_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/mapping_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/tokenizer_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/services"
	_ "github.com/NeF2le/anonix/gateway/swagger"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

// @title Anonix Gateway API
// @version 1.0
// @description Описание эндпоинтов сервиса анонимизации данных
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = logger.New(ctx)

	mainConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	tlsCfg := mainConfig.TLS

	ctx = context.WithValue(ctx, logger.KeyForLogLevel, mainConfig.LogLevel)

	tokenizerServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.Tokenizer, mainConfig.UpstreamPorts.Tokenizer)
	mappingServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.Mapping, mainConfig.UpstreamPorts.Mapping)
	authServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.Auth, mainConfig.UpstreamPorts.Auth)

	tokenizerServiceRepo := tokenizer_service_adapters.NewTokenizerServiceAdapterGRPC(
		tokenizerServiceAddress,
		mainConfig.Timeouts.Tokenizer,
	)
	mappingServiceRepo := mapping_service_adapters.NewMappingServiceAdapterGRPC(
		mappingServiceAddress,
		mainConfig.Timeouts.Mapping,
	)
	authServiceRepo := auth_service_adapters.NewAuthServiceAdapterGRPC(
		authServiceAddress,
		mainConfig.Timeouts.Auth,
	)

	if tlsCfg.Enabled {
		if err = tls_helpers.Verification(mainConfig.Gateway.Host, &tlsCfg); err != nil {
			panic(err)
		}
		var tlsClientCreds credentials.TransportCredentials
		tlsClientCreds, err = tls_helpers.LoadClientTLSConfig(tlsCfg.ClientPublicKey, tlsCfg.ClientPrivateKey, tlsCfg.RootPublicKey)
		if err != nil {
			panic(err)
		}

		tokenizerServiceRepo.AddTLS(tlsClientCreds)
		mappingServiceRepo.AddTLS(tlsClientCreds)
		authServiceRepo.AddTLS(tlsClientCreds)
	}

	tokenizerService := services.NewTokenizerService(
		tokenizerServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		mainConfig.GrpcPool.BaseRetryDelayMilliseconds,
	)
	mappingService := services.NewMappingService(
		mappingServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		mainConfig.GrpcPool.BaseRetryDelayMilliseconds,
	)
	authService := services.NewAuthService(
		authServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		mainConfig.GrpcPool.BaseRetryDelayMilliseconds,
	)

	tokenizerServiceHandler := http_handlers.NewTokenizerServiceHandler(tokenizerService, mappingService)
	mappingServiceHandler := http_handlers.NewMappingServiceHandler(mappingService)
	authServiceHandler := http_handlers.NewAuthServiceHandler(authService)

	authMiddleware := middlewares.NewAuthMiddleware(
		mainConfig.JWTSecret,
		authService,
		mainConfig.AccessTokenCookieTTL,
		mainConfig.RefreshTokenCookieTTL,
	)

	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	adminGroup := app.Group("/admin")
	adminGroup.Static("/", "static")
	adminGroup.GET("/scripts/config.js", func(c echo.Context) error {
		jsConfig := map[string]string{
			"API_BASE_URL": fmt.Sprintf("https://%s:%d/api/v1", mainConfig.Gateway.Host, mainConfig.HTTPPort),
			"APP_ENV":      mainConfig.Mode,
		}

		js := "window.APP_CONFIG = {\n"
		for key, value := range jsConfig {
			js += fmt.Sprintf("    %s: %q,\n", key, value)
		}
		js += "};\n"

		return c.Blob(http.StatusOK, "application/javascript", []byte(js))
	})

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	v1Group.Use(middlewares.LoggingMiddleware)

	tokenizerGroup := v1Group.Group("/tokenizer")
	tokenizerGroup.Use(authMiddleware.CheckAuth)
	{
		tokenizerGroup.POST("/tokenize", tokenizerServiceHandler.Tokenize)
		tokenizerGroup.POST("/detokenize", tokenizerServiceHandler.Detokenize)
	}

	mappingGroup := v1Group.Group("/mappings")
	mappingGroup.Use(authMiddleware.CheckAuth)
	{
		mappingGroup.GET("/:id", mappingServiceHandler.GetMapping)
		mappingGroup.GET("/", mappingServiceHandler.GetMappingList)
		mappingGroup.DELETE("/:id", mappingServiceHandler.DeleteMapping)
		mappingGroup.PATCH("/:id", mappingServiceHandler.UpdateMapping)
	}

	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/signIn", authServiceHandler.Login)
		authGroup.POST("/signUp", authServiceHandler.Register)
		authGroup.POST("/refresh", authServiceHandler.Refresh)
	}

	userGroup := v1Group.Group("/user")
	{
		userGroup.POST("/isAdmin", authServiceHandler.IsAdmin)
	}

	if tlsCfg.Enabled {
		rootCertFile := tlsCfg.RootPublicKey
		rootKeyFile := tlsCfg.RootPrivateKey

		if rootKeyFile == "" || rootCertFile == "" {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx,
				"TLS mode enabled but TLS_CERT_FILE or TLS_KEY_FILE is not set",
				slog.String("TLS_KEY_FILE", rootKeyFile),
				slog.String("TLS_CERT_FILE", rootCertFile))
		}

		app.Server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}

		go func() {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "server starting (TLS)",
				slog.Int("port", mainConfig.HTTPPort))

			if err = app.StartTLS(fmt.Sprintf(":%d", mainConfig.HTTPPort), rootCertFile, rootKeyFile); err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx, "error starting app (TLS)", logger.Err(err))
			}
		}()
	} else {
		go func() {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "server starting", slog.Int("port", mainConfig.HTTPPort))

			if err = app.Start(fmt.Sprintf(":%d", mainConfig.HTTPPort)); err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx, "error starting app", logger.Err(err))
			}
		}()
	}

	<-ctx.Done()

	err = app.Shutdown(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "error shutting down", logger.Err(err))
	}
}
