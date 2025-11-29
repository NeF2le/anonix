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

// UpdateMapping godoc
// @Summary Обновить маппинг по ID
// @Description Возвращает объект маппинга
// @Tags Mappings
// @Produce json
// @Param id path string true "ID маппинга"
// @Param body body schemas.UpdateMappingSchema true "Данные для обновления"
// @Success 200 {object} schemas.MappingSchema
// @Failure 400 "invalid token ID"
// @Failure 404 "mapping not found"
// @Failure 500 "internal error"
// @Security ApiKeyAuth
// @Router /mappings/{id} [patch]
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

	resp, err := m.mappingService.UpdateMapping(reqCtx, updateMappingReq)
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

	return ctx.JSON(http.StatusOK, helpers.ProtoMappingToSchema(resp.MappingModel))
}

// DeleteMapping godoc
// @Summary Удалить маппинг по ID
// @Description Удаляет объект маппинга по ID.
// @Tags Mappings
// @Produce json
// @Param id path string true "ID маппинга"
// @Success 200 {string} string "OK"
// @Failure 400 "invalid token ID"
// @Failure 401 "unauthorized"
// @Failure 404 "mapping not found"
// @Failure 500 "internal error"
// @Security ApiKeyAuth
// @Router /mappings/{id} [delete]
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

// GetMapping godoc
// @Summary Получить маппинг по ID
// @Description Возвращает объект маппинга
// @Tags Mappings
// @Produce json
// @Param id path string true "ID маппинга"
// @Success 200 {object} schemas.MappingSchema
// @Failure 400 "invalid token ID"
// @Failure 404 "mapping not found"
// @Failure 500 "internal error"
// @Security ApiKeyAuth
// @Router /mappings/{id} [get]
func (m *MappingServiceHandler) GetMapping(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	id, err := helpers.ParseUUID(ctx.Param("id"))
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to parse token UUID",
			slog.String("mapping ID", id.String()),
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid token ID")
	}

	resp, err := m.mappingService.GetMapping(reqCtx, &mapping.GetMappingRequest{Id: id.String()})
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

	return ctx.JSON(http.StatusOK, helpers.ProtoMappingToSchema(resp.MappingModel))
}

// GetMappingList godoc
// @Summary Получить список маппингов
// @Description Возвращает список маппингов
// @Tags Mappings
// @Produce json
// @Success 200 {array} schemas.MappingSchema
// @Failure 500 "failed to get mapping list"
// @Security ApiKeyAuth
// @Router /mappings/ [get]
func (m *MappingServiceHandler) GetMappingList(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	resp, err := m.mappingService.GetMappingList(reqCtx, &mapping.GetMappingListRequest{})
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to get mapping list", logger.Err(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get mapping list")
	}

	var mappings []*schemas.MappingSchema
	for _, mm := range resp.MappingModels {
		mappings = append(mappings, helpers.ProtoMappingToSchema(mm))
	}
	return ctx.JSON(http.StatusOK, mappings)
}
