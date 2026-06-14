package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/NeF2le/anonix/gateway/internal/config"
	"github.com/NeF2le/anonix/gateway/internal/domain"
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
	keyRotationHandler := http_handlers.NewKeyRotationHandler(tokenizerService, mappingService)

	authMiddleware := middlewares.NewAuthMiddleware(
		mainConfig.JWTSecret,
		authService,
		mainConfig.AccessTokenCookieTTL,
		mainConfig.RefreshTokenCookieTTL,
	)
	rbacMiddleware := middlewares.NewRBACMiddleware()

	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	adminGroup := app.Group("/admin")
	adminGroup.Static("/", "static")
	var apiBaseURL string
	if tlsCfg.Enabled {
		apiBaseURL = fmt.Sprintf("https://%s:%d/api/v1", mainConfig.Gateway.Host, mainConfig.HTTPPort)
	} else {
		apiBaseURL = fmt.Sprintf("http://%s:%d/api/v1", mainConfig.Gateway.Host, mainConfig.HTTPPort)
	}
	adminGroup.GET("/scripts/config.js", func(c echo.Context) error {
		jsConfig := map[string]string{
			"API_BASE_URL": apiBaseURL,
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
	tokenizerGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin, domain.RoleSpecialist))
	{
		tokenizerGroup.POST("/tokenize", tokenizerServiceHandler.Tokenize)
		tokenizerGroup.POST("/detokenize", tokenizerServiceHandler.Detokenize)
	}

	mappingReadGroup := v1Group.Group("/mappings")
	mappingReadGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin, domain.RoleSpecialist, domain.RoleAuditor))
	{
		mappingReadGroup.GET("/:id", mappingServiceHandler.GetMapping)
		mappingReadGroup.GET("/", mappingServiceHandler.GetMappingList)
	}

	mappingWriteGroup := v1Group.Group("/mappings")
	mappingWriteGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin, domain.RoleSpecialist))
	{
		mappingWriteGroup.DELETE("/:id", mappingServiceHandler.DeleteMapping)
		mappingWriteGroup.PATCH("/:id", mappingServiceHandler.UpdateMapping)
	}

	kindReadGroup := v1Group.Group("/kinds")
	kindReadGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin, domain.RoleSpecialist, domain.RoleAuditor))
	{
		kindReadGroup.GET("/:id", mappingServiceHandler.GetKind)
		kindReadGroup.GET("/", mappingServiceHandler.GetKindList)
	}

	kindWriteGroup := v1Group.Group("/kinds")
	kindWriteGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin))
	{
		kindWriteGroup.POST("/", mappingServiceHandler.CreateKind)
		kindWriteGroup.PATCH("/:id", mappingServiceHandler.UpdateKind)
		kindWriteGroup.DELETE("/:id", mappingServiceHandler.DeleteKind)
	}

	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/signIn", authServiceHandler.Login)
		authGroup.POST("/signUp", authServiceHandler.Register)
		authGroup.POST("/refresh", authServiceHandler.Refresh)
		authGroup.GET("/me", authServiceHandler.GetMe, authMiddleware.CheckAuth)
	}

	userGroup := v1Group.Group("/user")
	userGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin))
	{
		userGroup.POST("/isAdmin", authServiceHandler.IsAdmin)
		userGroup.GET("/list", authServiceHandler.GetUsers)
		userGroup.DELETE("/delete", authServiceHandler.DeleteUser)
		userGroup.GET("/roles", authServiceHandler.GetUserRoles)
		userGroup.POST("/assignRole", authServiceHandler.AssignRole)
		userGroup.DELETE("/removeRole", authServiceHandler.RemoveRole)
		userGroup.PATCH("/clearance", authServiceHandler.UpdateClearance)
	}

	roleGroup := v1Group.Group("/role")
	roleGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin))
	{
		roleGroup.GET("/list", authServiceHandler.GetRolesList)
	}

	auditGroup := v1Group.Group("/audit")
	auditGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin, domain.RoleAuditor))
	{
		auditGroup.GET("/", mappingServiceHandler.GetAuditLogList)
	}

	keysGroup := v1Group.Group("/admin/keys")
	keysGroup.Use(authMiddleware.CheckAuth, rbacMiddleware.CheckRole(domain.RoleAdmin))
	{
		keysGroup.POST("/rotate-master", keyRotationHandler.RotateMasterKey)
		keysGroup.POST("/rotate-deks", keyRotationHandler.RotateAllDeks)
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
