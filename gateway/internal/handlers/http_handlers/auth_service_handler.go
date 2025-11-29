package http_handlers

import (
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/handlers/helpers"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

type AuthServiceHandler struct {
	authService        *services.AuthService
	accessTokenMaxAge  int
	refreshTokenMaxAge int
}

func NewAuthServiceHandler(authService *services.AuthService) *AuthServiceHandler {
	return &AuthServiceHandler{authService: authService}
}

// Register godoc
// @Summary Зарегистрировать пользователя
// @Description Возвращает ID пользователя
// @Tags Auth
// @Produce json
// @Param body body schemas.RegisterSchema true "Данные для регистрации"
// @Success 200 {object} schemas.RegisterRespSchema
// @Failure 400 "invalid request body"
// @Failure 409 "user with same login already exists"
// @Failure 500 "failed to register user"
// @Router /auth/signUp [post]
func (a *AuthServiceHandler) Register(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var registerSchema *schemas.RegisterSchema
	err := ctx.Bind(&registerSchema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind RegisterSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	registerReq := &auth_service.RegisterRequest{
		Login:    registerSchema.Login,
		Password: registerSchema.Password,
		RoleId:   int32(registerSchema.RoleId),
	}

	registerResp, err := a.authService.Register(reqCtx, registerReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("login", registerReq.Login),
					slog.String("password", registerReq.Password),
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.AlreadyExists:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "already registered",
					slog.String("login", registerReq.Login),
					slog.String("password", registerReq.Password),
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.Conflict(ctx, "user with same login already exists")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to register user",
					slog.String("login", registerReq.Login),
					slog.String("password", registerReq.Password),
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to register user")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while registering user",
			slog.String("login", registerReq.Login),
			slog.String("password", registerReq.Password),
			slog.Int("roleId", int(registerReq.RoleId)),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.RegisterRespSchema{UserId: registerResp.UserId})
}

// Login godoc
// @Summary Авторизация
// @Description Возвращает токен доступа и токен обновления
// @Tags Auth
// @Produce json
// @Param body body schemas.LoginSchema true "Данные для авторизации"
// @Success 200 {object} schemas.LoginRespSchema
// @Failure 400 "invalid request body"
// @Failure 401 "invalid credentials"
// @Failure 500 "failed to log in"
// @Router /auth/signIn [post]
func (a *AuthServiceHandler) Login(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var loginSchema *schemas.LoginSchema
	err := ctx.Bind(&loginSchema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind LoginSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	loginReq := &auth_service.LoginRequest{
		Login:    loginSchema.Login,
		Password: loginSchema.Password,
	}

	loginResp, err := a.authService.Login(reqCtx, loginReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("login", loginReq.Login),
					slog.String("password", loginReq.Password),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.Unauthenticated:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "unauthenticated user",
					slog.String("login", loginReq.Login),
					slog.String("password", loginReq.Password),
					logger.Err(err))
				return helpers.InvalidCredentials(ctx)
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to log in",
					slog.String("login", loginReq.Login),
					slog.String("password", loginReq.Password),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to log in")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while logging user",
			slog.String("login", loginReq.Login),
			slog.String("password", loginReq.Password),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	helpers.SetAccessTokenCookie(ctx, loginResp.AccessToken, a.accessTokenMaxAge)
	helpers.SetRefreshTokenCookie(ctx, loginResp.RefreshToken, a.refreshTokenMaxAge)

	return ctx.JSON(http.StatusOK, &schemas.LoginRespSchema{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		UserId:       loginResp.UserId,
	})
}

// Refresh godoc
// @Summary Обновить токены
// @Description Принимает refresh token, возвращает новый access и refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body schemas.RefreshSchema true "Refresh token"
// @Success 200 {object} schemas.RefreshRespSchema
// @Failure 400 "invalid request body / failed to refresh token"
// @Failure 500 "unexpected error"
// @Router /auth/refresh [post]
func (a *AuthServiceHandler) Refresh(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var refreshSchema *schemas.RefreshSchema
	err := ctx.Bind(&refreshSchema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind RefreshSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	refreshReq := &auth_service.RefreshRequest{RefreshToken: refreshSchema.RefreshToken}

	refreshResp, err := a.authService.Refresh(reqCtx, refreshReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("refresh", refreshReq.RefreshToken),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to refresh token",
					slog.String("refresh", refreshReq.RefreshToken),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to refresh token")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while refreshing auth token",
			slog.String("refresh", refreshReq.RefreshToken),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.RefreshRespSchema{
		AccessToken:  refreshResp.AccessToken,
		RefreshToken: refreshResp.RefreshToken,
	})
}

// IsAdmin godoc
// @Summary Проверить админство пользователя
// @Description По переданному user_id возвращает, является ли пользователь администратором
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body schemas.IsAdminSchema true "ID пользователя"
// @Success 200 {object} schemas.IsAdminRespSchema
// @Failure 400 "invalid request body"
// @Failure 500 "failed to check if user is admin"
// @Router /auth/isAdmin [post]
func (a *AuthServiceHandler) IsAdmin(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var isAdminSchema *schemas.IsAdminSchema
	err := ctx.Bind(&isAdminSchema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind IsAdminSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	isAdminReq := &auth_service.IsAdminRequest{UserId: isAdminSchema.UserId}

	isAdminResp, err := a.authService.IsAdmin(reqCtx, isAdminReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", isAdminReq.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to check if user is admin",
					slog.String("user ID", isAdminReq.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to check if user is admin")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while checking if user is admin",
			slog.String("user ID", isAdminReq.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.IsAdminRespSchema{Result: isAdminResp.Result})
}
