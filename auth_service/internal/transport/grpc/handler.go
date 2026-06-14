package grpc

import (
	"context"
	"errors"
	"github.com/NeF2le/anonix/auth_service/internal/ports"
	"github.com/NeF2le/anonix/auth_service/internal/transport/helpers"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"
)

type grpcAuthHandler struct {
	auth ports.AuthUseCase
	auth_service.UnimplementedAuthServiceServer
}

func NewGRPCAuthHandler(useCase ports.AuthUseCase) auth_service.AuthServiceServer {
	return &grpcAuthHandler{auth: useCase}
}

func (s *grpcAuthHandler) Register(ctx context.Context, req *auth_service.RegisterRequest) (
	*auth_service.RegisterResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login required")
	}
	if strings.ContainsRune(req.GetLogin(), ' ') {
		return nil, status.Error(codes.InvalidArgument, "login contains space")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}
	if strings.ContainsRune(req.GetPassword(), ' ') {
		return nil, status.Error(codes.InvalidArgument, "password contains space")
	}
	if req.GetRoleId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "role required and must be greater than zero")
	}
	if err := helpers.ValidatePassword(req.GetPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"user already exists",
				slog.String("login", req.GetLogin()),
				logger.Err(err),
			)
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to register user",
			slog.String("login", req.GetLogin()),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	if err = s.auth.AssignRole(ctx, userID, int(req.GetRoleId())); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to assign role after registration",
			slog.String("userID", userID),
			slog.Int("roleId", int(req.GetRoleId())),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to assign role")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"user registered successfully",
		slog.String("login", req.GetLogin()),
		slog.String("userID", userID),
	)
	return &auth_service.RegisterResponse{UserId: userID}, nil
}

func (s *grpcAuthHandler) Login(ctx context.Context, req *auth_service.LoginRequest) (
	*auth_service.LoginResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	userID, accessToken, refreshToken, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"invalid credentials",
				slog.String("login", req.GetLogin()),
				logger.Err(err),
			)
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to login user",
			slog.String("login", req.GetLogin()),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to login")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"user logged in successfully",
		slog.String("login", req.GetLogin()),
		slog.String("userID", userID),
	)
	return &auth_service.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *grpcAuthHandler) Refresh(ctx context.Context, req *auth_service.RefreshRequest) (
	*auth_service.RefreshResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, errs.ErrInvalidToken
	}

	accessToken, refreshToken, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, errs.ErrInvalidToken) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "invalid token",
				logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, errs.ErrTokenExpired) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "token expired",
				logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to refresh token",
			logger.Err(err))
		return nil, status.Error(codes.Unauthenticated, "failed to refresh token")
	}

	return &auth_service.RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *grpcAuthHandler) IsAdmin(ctx context.Context, req *auth_service.IsAdminRequest) (
	*auth_service.IsAdminResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "user doesn't exists",
				slog.String("userId", req.GetUserId()),
				logger.Err(err))
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, errs.ErrInvalidCredentials) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "invalid user ID",
				slog.String("userId", req.GetUserId()),
				logger.Err(err))
			return nil, status.Error(codes.InvalidArgument, "invalid user ID")
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to check if user is admin",
			slog.String("userId", req.GetUserId()),
			logger.Err(err))
		return nil, status.Error(codes.Internal, "failed to check user admin")
	}

	return &auth_service.IsAdminResponse{Result: isAdmin}, nil
}

func (s *grpcAuthHandler) GetUsers(ctx context.Context, req *auth_service.GetUsersRequest) (
	*auth_service.GetUsersResponse, error) {

	users, err := s.auth.GetUsers(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get users", logger.Err(err))
		return nil, status.Error(codes.Internal, "failed to get users")
	}

	result := make([]*auth_service.User, 0, len(users))
	for _, u := range users {
		protoUser := &auth_service.User{
			Id:             u.ID,
			Login:          u.Login,
			ClearanceLevel: int32(u.ClearanceLevel),
		}
		for _, r := range u.Roles {
			protoUser.Roles = append(protoUser.Roles, &auth_service.Role{
				Id:   int32(r.ID),
				Name: r.Name,
			})
		}
		result = append(result, protoUser)
	}

	return &auth_service.GetUsersResponse{Users: result}, nil
}

func (s *grpcAuthHandler) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (
	*auth_service.DeleteUserResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}

	err := s.auth.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"user not found",
				slog.String("userId", req.GetUserId()),
				logger.Err(err),
			)
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to delete user",
			slog.String("userId", req.GetUserId()),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"user deleted successfully",
		slog.String("userId", req.GetUserId()),
	)

	return &auth_service.DeleteUserResponse{}, nil
}

func (s *grpcAuthHandler) AssignRole(ctx context.Context, req *auth_service.AssignRoleRequest) (
	*auth_service.AssignRoleResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}
	if req.GetRoleId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "role ID must be greater than zero")
	}

	err := s.auth.AssignRole(ctx, req.GetUserId(), int(req.GetRoleId()))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"user not found",
				slog.String("userId", req.GetUserId()),
				logger.Err(err),
			)
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to assign role",
			slog.String("userId", req.GetUserId()),
			slog.Int("roleId", int(req.GetRoleId())),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to assign role")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"role assigned successfully",
		slog.String("userId", req.GetUserId()),
		slog.Int("roleId", int(req.GetRoleId())),
	)

	return &auth_service.AssignRoleResponse{}, nil
}

func (s *grpcAuthHandler) UpdateClearanceLevel(ctx context.Context, req *auth_service.UpdateClearanceLevelRequest) (
	*auth_service.UpdateClearanceLevelResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}
	if req.GetClearanceLevel() < 1 || req.GetClearanceLevel() > 4 {
		return nil, status.Error(codes.InvalidArgument, "clearance level must be between 1 and 4")
	}

	err := s.auth.UpdateClearanceLevel(ctx, req.GetUserId(), int(req.GetClearanceLevel()))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"user not found",
				slog.String("userId", req.GetUserId()),
				logger.Err(err),
			)
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to update clearance level",
			slog.String("userId", req.GetUserId()),
			slog.Int("clearanceLevel", int(req.GetClearanceLevel())),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to update clearance level")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"clearance level updated successfully",
		slog.String("userId", req.GetUserId()),
		slog.Int("clearanceLevel", int(req.GetClearanceLevel())),
	)

	return &auth_service.UpdateClearanceLevelResponse{}, nil
}

func (s *grpcAuthHandler) RemoveRole(ctx context.Context, req *auth_service.RemoveRoleRequest) (
	*auth_service.RemoveRoleResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID required")
	}
	if req.GetRoleId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "role ID must be greater than zero")
	}

	err := s.auth.RemoveRole(ctx, req.GetUserId(), int(req.GetRoleId()))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx,
				"user not found",
				slog.String("userId", req.GetUserId()),
				logger.Err(err),
			)
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to remove role",
			slog.String("userId", req.GetUserId()),
			slog.Int("roleId", int(req.GetRoleId())),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to remove role")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"role removed successfully",
		slog.String("userId", req.GetUserId()),
		slog.Int("roleId", int(req.GetRoleId())),
	)

	return &auth_service.RemoveRoleResponse{}, nil
}

func (s *grpcAuthHandler) GetRolesList(ctx context.Context, req *auth_service.GetRolesListRequest) (
	*auth_service.GetRolesListResponse, error) {

	roles, err := s.auth.GetRolesList(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get roles list",
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to get roles list")
	}

	result := make([]*auth_service.Role, 0, len(roles))
	for _, r := range roles {
		result = append(result, &auth_service.Role{
			Id:   int32(r.ID),
			Name: r.Name,
		})
	}

	return &auth_service.GetRolesListResponse{
		Roles: result,
	}, nil
}

func (s *grpcAuthHandler) GetUserRoles(ctx context.Context, req *auth_service.GetUserRolesRequest) (
	*auth_service.GetUserRolesResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}

	roles, err := s.auth.GetUserRoles(ctx, req.GetUserId())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get user roles",
			slog.String("userId", req.GetUserId()),
			logger.Err(err),
		)
		return nil, status.Error(codes.Internal, "failed to get user roles")
	}

	result := make([]*auth_service.Role, 0, len(roles))
	for _, r := range roles {
		result = append(result, &auth_service.Role{
			Id:   int32(r.ID),
			Name: r.Name,
		})
	}

	return &auth_service.GetUserRolesResponse{
		Roles: result,
	}, nil
}
