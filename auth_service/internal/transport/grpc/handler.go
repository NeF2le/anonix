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

	userID, err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword(), int(req.GetRoleId()))
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
		slog.String("accessToken", accessToken),
		slog.String("refreshToken", refreshToken),
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
				slog.String("token", req.GetRefreshToken()),
				logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, errs.ErrTokenExpired) {
			logger.GetLoggerFromCtx(ctx).Warn(ctx, "token expired",
				slog.String("token", req.GetRefreshToken()),
				logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to refresh token",
			slog.String("refreshToken", req.GetRefreshToken()),
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
