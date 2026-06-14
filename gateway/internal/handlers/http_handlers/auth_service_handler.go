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
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.AlreadyExists:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "already registered",
					slog.String("login", registerReq.Login),
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.Conflict(ctx, "user with same login already exists")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to register user",
					slog.String("login", registerReq.Login),
					slog.Int("roleId", int(registerReq.RoleId)),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to register user")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while registering user",
			slog.String("login", registerReq.Login),
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
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.Unauthenticated:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "unauthenticated user",
					slog.String("login", loginReq.Login),
					logger.Err(err))
				return helpers.InvalidCredentials(ctx)
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to log in",
					slog.String("login", loginReq.Login),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to log in")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while logging user",
			slog.String("login", loginReq.Login),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	helpers.SetAccessTokenCookie(ctx, loginResp.AccessToken, a.accessTokenMaxAge)
	helpers.SetRefreshTokenCookie(ctx, loginResp.RefreshToken, a.refreshTokenMaxAge)

	rolesResp, err := a.authService.GetUserRoles(reqCtx, &auth_service.GetUserRolesRequest{UserId: loginResp.UserId})
	if err != nil {
		rolesResp = &auth_service.GetUserRolesResponse{}
	}

	roleSchemas := make([]*schemas.RoleSchema, 0, len(rolesResp.Roles))
	for _, r := range rolesResp.Roles {
		roleSchemas = append(roleSchemas, &schemas.RoleSchema{Id: r.Id, Name: r.Name})
	}

	return ctx.JSON(http.StatusOK, &schemas.LoginRespSchema{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		UserId:       loginResp.UserId,
		Roles:        roleSchemas,
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
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to refresh token",
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to refresh token")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while refreshing auth token",
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

// GetMe godoc
// @Summary Получить роли текущего пользователя
// @Description Возвращает роли авторизованного пользователя из токена доступа
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.GetUserRolesRespSchema
// @Failure 401 "unauthorized"
// @Router /auth/me [get]
func (a *AuthServiceHandler) GetMe(ctx echo.Context) error {
	protoRoles, _ := ctx.Get("roles").([]*auth_service.Role)

	roles := make([]*schemas.RoleSchema, 0, len(protoRoles))
	for _, r := range protoRoles {
		roles = append(roles, helpers.ProtoRoleToSchema(r))
	}

	return ctx.JSON(http.StatusOK, &schemas.GetUserRolesRespSchema{Roles: roles})
}

func (a *AuthServiceHandler) GetUsers(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	resp, err := a.authService.GetUsers(reqCtx, &auth_service.GetUsersRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to get users",
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to get users")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while getting users",
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	users := make([]*schemas.UserSchema, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, helpers.ProtoUserToSchema(u))
	}

	return ctx.JSON(http.StatusOK, &schemas.GetUsersRespSchema{Users: users})
}

func (a *AuthServiceHandler) DeleteUser(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var schema *schemas.DeleteUserSchema
	err := ctx.Bind(&schema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind DeleteUserSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	req := &auth_service.DeleteUserRequest{UserId: schema.UserId}

	_, err = a.authService.DeleteUser(reqCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "user not found",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "user not found")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to delete user",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to delete user")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while deleting user",
			slog.String("user ID", req.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.DeleteUserRespSchema{})
}

func (a *AuthServiceHandler) AssignRole(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var schema *schemas.AssignRoleSchema
	err := ctx.Bind(&schema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind AssignRoleSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	req := &auth_service.AssignRoleRequest{UserId: schema.UserId, RoleId: schema.RoleId}

	_, err = a.authService.AssignRole(reqCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "user not found",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "user not found")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to assign role",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to assign role")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while assigning role",
			slog.String("user ID", req.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.AssignRoleRespSchema{})
}

func (a *AuthServiceHandler) RemoveRole(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var schema *schemas.RemoveRoleSchema
	err := ctx.Bind(&schema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind RemoveRoleSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	req := &auth_service.RemoveRoleRequest{UserId: schema.UserId, RoleId: schema.RoleId}

	_, err = a.authService.RemoveRole(reqCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "user not found",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "user not found")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to remove role",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to remove role")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while removing role",
			slog.String("user ID", req.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.RemoveRoleRespSchema{})
}

func (a *AuthServiceHandler) UpdateClearance(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var schema *schemas.UpdateClearanceSchema
	err := ctx.Bind(&schema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind UpdateClearanceSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	if schema.ClearanceLevel < 1 || schema.ClearanceLevel > 4 {
		return helpers.BadRequest(ctx, "invalid request body")
	}

	req := &auth_service.UpdateClearanceLevelRequest{UserId: schema.UserId, ClearanceLevel: schema.ClearanceLevel}

	_, err = a.authService.UpdateClearanceLevel(reqCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "user not found",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "user not found")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to update clearance level",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to update clearance level")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while updating clearance level",
			slog.String("user ID", req.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.UpdateClearanceRespSchema{})
}

func (a *AuthServiceHandler) GetRolesList(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	resp, err := a.authService.GetRolesList(reqCtx, &auth_service.GetRolesListRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to get roles list",
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to get roles list")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while getting roles list",
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	roles := make([]*schemas.RoleSchema, 0, len(resp.Roles))
	for _, r := range resp.Roles {
		roles = append(roles, helpers.ProtoRoleToSchema(r))
	}

	return ctx.JSON(http.StatusOK, &schemas.GetRolesListRespSchema{Roles: roles})
}

func (a *AuthServiceHandler) GetUserRoles(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var schema *schemas.GetUserRolesSchema
	err := ctx.Bind(&schema)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind GetUserRolesSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	req := &auth_service.GetUserRolesRequest{UserId: schema.UserId}

	resp, err := a.authService.GetUserRoles(reqCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid request body",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "user not found",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "user not found")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to get user roles",
					slog.String("user ID", req.UserId),
					logger.Err(err))
				return helpers.BadRequest(ctx, "failed to get user roles")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while getting user roles",
			slog.String("user ID", req.UserId),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	roles := make([]*schemas.RoleSchema, 0, len(resp.Roles))
	for _, r := range resp.Roles {
		roles = append(roles, helpers.ProtoRoleToSchema(r))
	}

	return ctx.JSON(http.StatusOK, &schemas.GetUserRolesRespSchema{Roles: roles})
}
