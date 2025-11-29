package http_handlers

import (
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
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
	"time"
)

type TokenizerServiceHandler struct {
	tokenizerService *services.TokenizerService
	mappingService   *services.MappingService
}

func NewTokenizerServiceHandler(
	tokenizerService *services.TokenizerService,
	mappingService *services.MappingService) *TokenizerServiceHandler {
	return &TokenizerServiceHandler{
		tokenizerService: tokenizerService,
		mappingService:   mappingService,
	}
}

// Tokenize godoc
// @Summary Токенизация
// @Description Принимает plaintext и параметры токенизации, возвращает созданный токен и метаданные
// @Tags Tokenizer
// @Accept json
// @Produce json
// @Param body body schemas.TokenizeSchema true "Данные для токенизации"
// @Success 200 {object} schemas.MappingSchema
// @Failure 400 "invalid request body / invalid arguments"
// @Failure 409 "token already exists"
// @Failure 500 "failed to tokenize / unexpected error"
// @Security ApiKeyAuth
// @Router /tokenize [post]
func (t *TokenizerServiceHandler) Tokenize(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var tokenizeSchema *schemas.TokenizeSchema
	if err := ctx.Bind(&tokenizeSchema); err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to bind tokenize schema",
			logger.Err(err))
		return helpers.BadRequest(ctx, "invalid request body")
	}

	tokenizeReq := &tokenizer.TokenizeRequest{
		Plaintext:     tokenizeSchema.Plaintext,
		Deterministic: tokenizeSchema.Deterministic,
		Reversible:    tokenizeSchema.Reversible,
	}

	tokenizeResp, err := t.tokenizerService.Tokenize(reqCtx, tokenizeReq)
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "tokenizerService.Tokenize failed",
			logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to tokenize")
	}

	mappingReq := &mapping.CreateMappingRequest{
		CipherText:    tokenizeResp.CipherText,
		DekWrapped:    tokenizeResp.DekWrapped,
		Reversible:    tokenizeResp.Reversible,
		Deterministic: tokenizeResp.Deterministic,
		TokenTtl:      durationpb.New(time.Duration(tokenizeSchema.TokenTTL) * time.Second),
	}
	resp, err := t.mappingService.CreateMapping(reqCtx, mappingReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "token already exists", logger.Err(err))
				return helpers.Conflict(ctx, "token already exists")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "invalid arguments", logger.Err(err))
				return helpers.BadRequest(ctx, "invalid arguments")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to create mapping",
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to tokenize")
			}
		}
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type", logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, helpers.ProtoMappingToSchema(resp.MappingModel))
}

// Detokenize godoc
// @Summary Детокенизация
// @Description Принимает токен, ищет соответствующий mapping и возвращает исходный plaintext
// @Tags Tokenizer
// @Accept json
// @Produce json
// @Param body body schemas.DetokenizeSchema true "Токен для детокенизации"
// @Success 200 {object} schemas.DetokenizeRespSchema
// @Failure 400 "invalid request body / invalid arguments"
// @Failure 404 "token not found / token expired"
// @Failure 500 "failed to detokenize / unexpected error"
// @Security ApiKeyAuth
// @Router /detokenize [post]
func (t *TokenizerServiceHandler) Detokenize(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	var detokenizeSchema *schemas.DetokenizeSchema
	if err := ctx.Bind(&detokenizeSchema); err != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx,
			"failed to bind detokenize schema",
			logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to detokenize")
	}

	getMappingReq := &mapping.GetMappingRequest{Id: detokenizeSchema.Token}
	getMappingResp, err := t.mappingService.GetMapping(reqCtx, getMappingReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping not found",
					slog.String("mapping ID", detokenizeSchema.Token),
					logger.Err(err))
				return helpers.NotFound(ctx, "token not found")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "invalid mapping ID",
					slog.String("mapping ID", detokenizeSchema.Token),
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid arguments")
			case codes.DeadlineExceeded:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "token expired",
					slog.String("mapping ID", detokenizeSchema.Token),
					logger.Err(err))
				return helpers.NotFound(ctx, "token expired")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to detokenize", logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to detokenize")
			}
		}
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type", logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	detokenizeReq := &tokenizer.DetokenizeRequest{
		CipherText:    getMappingResp.MappingModel.CipherText,
		DekWrapped:    getMappingResp.MappingModel.DekWrapped,
		Deterministic: getMappingResp.MappingModel.Deterministic,
	}

	detokenizeResp, err := t.tokenizerService.Detokenize(reqCtx, detokenizeReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "failed to detokenize",
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid request body")
			default:
				logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "tokenizerService.Detokenize failed",
					logger.Err(err))
				return helpers.InternalServerError(ctx, "failed to detokenize")
			}
		}
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "unexpected error type", logger.Err(err))
		return helpers.InternalServerError(ctx, "unexpected error")
	}

	return ctx.JSON(http.StatusOK, &schemas.DetokenizeRespSchema{Plaintext: detokenizeResp.Plaintext})
}
