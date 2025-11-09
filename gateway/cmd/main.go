package main

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/config"
	"github.com/NeF2le/anonix/gateway/internal/handlers/http_handlers"
	"github.com/NeF2le/anonix/gateway/internal/handlers/middlewares"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/auth_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/mapping_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/ports/adapters/tokenizer_service_adapters"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = logger.New(ctx)

	mainConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

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

	adminGroup := app.Group("/admin")
	adminGroup.Static("/", "static")
	adminGroup.GET("/scripts/config.js", func(c echo.Context) error {
		jsConfig := map[string]string{
			"API_BASE_URL": fmt.Sprintf("http://%s:%d/api/v1", mainConfig.Gateway.Host, mainConfig.HTTPPort),
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

	go func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "server starting", slog.Int("port", mainConfig.HTTPPort))

		//err = app.Start(fmt.Sprintf("%s:%d", mainConfig.Gateway.Host, mainConfig.HTTPPort))
		err = app.Start(fmt.Sprintf(":%d", mainConfig.HTTPPort))

		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "error starting app", logger.Err(err))
		}
	}()

	<-ctx.Done()

	err = app.Shutdown(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "error shutting down", logger.Err(err))
	}
}
