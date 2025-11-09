package http_handlers

import (
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/handlers/helpers"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"log/slog"
	"net/http"
)

type MappingServiceHandler struct {
	mappingService *services.MappingService
}

func NewMappingServiceHandler(mappingService *services.MappingService) *MappingServiceHandler {
	return &MappingServiceHandler{mappingService: mappingService}
}

func (m *MappingServiceHandler) UpdateMapping(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	id, err := helpers.ParseUUID(ctx.Param("id"))
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to parse token UUID",
			slog.String("mapping ID", id.String()),
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid token ID")
	}

	var updateMappingSchema *schemas.UpdateMappingSchema
	if err = ctx.Bind(&updateMappingSchema); err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind UpdateMappingSchema request body",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	updateMappingReq := &mapping.UpdateMappingRequest{
		Id:       id.String(),
		TokenTtl: durationpb.New(updateMappingSchema.TokenTtl),
	}

	updateMappingResp, err := m.mappingService.UpdateMapping(reqCtx, updateMappingReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping not found",
					slog.String("ID", id.String()))
				return helpers.NotFound(ctx, "mapping not found")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid arguments for update mapping",
					slog.String("ID", id.String()))
				return helpers.BadRequest(ctx, "invalid arguments for update mapping")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to update mapping",
					slog.String("ID", id.String()),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to update mapping")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type while updating mapping",
			slog.String("ID", id.String()),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, updateMappingResp)
}

func (m *MappingServiceHandler) DeleteMapping(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	id, err := helpers.ParseUUID(ctx.Param("id"))
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to parse token UUID",
			slog.String("mapping ID", id.String()),
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid token ID")
	}

	_, err = m.mappingService.DeleteMapping(reqCtx, &mapping.DeleteMappingRequest{Id: id.String()})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping not found",
					slog.String("ID", id.String()))
				return helpers.NotFound(ctx, "mapping not found")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid arguments for delete mapping",
					slog.String("ID", id.String()))
				return helpers.BadRequest(ctx, "invalid arguments for get mapping")
			default:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to delete mapping",
					slog.String("ID", id.String()),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to delete mapping")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type",
			slog.String("ID", id.String()),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (m *MappingServiceHandler) GetMapping(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	id, err := helpers.ParseUUID(ctx.Param("id"))
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to parse token UUID",
			slog.String("mapping ID", id.String()),
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid token ID")
	}

	getMappingResp, err := m.mappingService.GetMapping(reqCtx, &mapping.GetMappingRequest{Id: id.String()})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping not found",
					slog.String("ID", id.String()))
				return helpers.NotFound(ctx, "mapping not found")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "invalid arguments for get mapping",
					slog.String("ID", id.String()))
				return helpers.BadRequest(ctx, "invalid arguments for get mapping")
			case codes.DeadlineExceeded:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping expired",
					slog.String("ID", id.String()),
					logger.Err(err))
				return helpers.NotFound(ctx, "mapping expired")
			default:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to get mapping",
					slog.String("ID", id.String()),
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to get mapping")
			}
		}

		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type",
			slog.String("ID", id.String()),
			logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, getMappingResp.MappingModel)
}

func (m *MappingServiceHandler) GetMappingList(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	resp, err := m.mappingService.GetMappingList(reqCtx, &mapping.GetMappingListRequest{})
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to get mapping list", logger.Err(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get mapping list")
	}

	return ctx.JSON(http.StatusOK, resp.MappingModels)
}
