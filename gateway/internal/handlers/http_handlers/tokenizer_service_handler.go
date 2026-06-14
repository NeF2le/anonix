package http_handlers

import (
	"encoding/hex"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/domain"
	"github.com/NeF2le/anonix/gateway/internal/handlers/helpers"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

const (
	modePseudonymize = "pseudonymize"
	modeAnonymize    = "anonymize"
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
// @Description Принимает plaintext и параметры токенизации, возвращает созданный токен и метаданные.
// @Description Режим "pseudonymize" — обратимая операция, mapping сохраняется и его можно детокенизировать.
// @Description Режим "anonymize" — необратимая операция, mapping нигде не сохраняется,
// @Description ответом является schemas.TokenizeResultSchema, такой токен нельзя детокенизировать.
// @Tags Tokenizer
// @Accept json
// @Produce json
// @Param body body schemas.TokenizeSchema true "Данные для токенизации"
// @Success 200 {object} schemas.MappingSchema "mode=pseudonymize"
// @Success 200 {object} schemas.TokenizeResultSchema "mode=anonymize"
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

	pseudonymize := tokenizeSchema.Mode == modePseudonymize
	if !pseudonymize && tokenizeSchema.Mode != modeAnonymize {
		return helpers.BadRequest(ctx, "invalid mode")
	}

	switch tokenizeSchema.Algorithm {
	case "", "aes-siv", "gost-kuznechik":
	default:
		return helpers.BadRequest(ctx, "invalid algorithm")
	}

	var kind *mapping.Kind
	if tokenizeSchema.KindId > 0 {
		kindResp, err := t.mappingService.GetKind(reqCtx, &mapping.GetKindRequest{Id: int32(tokenizeSchema.KindId)})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				return helpers.BadRequest(ctx, "kind not found")
			}
			logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "mappingService.GetKind failed", logger.Err(err))
			return helpers.InternalServerError(ctx, "failed to tokenize")
		}
		kind = kindResp.Kind

		if !helpers.HasRole(ctx, domain.RoleAdmin) && helpers.GetClearanceLevel(ctx) < int(kind.AccessLevel) {
			return helpers.Forbidden(ctx, "insufficient clearance level")
		}

		if kind.Mask != "" {
			re, err := regexp.Compile(kind.Mask)
			if err != nil {
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "invalid kind mask regex",
					slog.String("mask", kind.Mask), logger.Err(err))
				return helpers.InternalServerError(ctx, "invalid kind mask")
			}
			if !re.Match(tokenizeSchema.Plaintext) {
				return helpers.BadRequest(ctx, "data does not match kind format")
			}
		}
	}

	tokenizeReq := &tokenizer.TokenizeRequest{
		Plaintext:     tokenizeSchema.Plaintext,
		Deterministic: tokenizeSchema.Deterministic,
		Pseudonymize:  pseudonymize,
		Algorithm:     tokenizeSchema.Algorithm,
	}

	tokenizeResp, err := t.tokenizerService.Tokenize(reqCtx, tokenizeReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "invalid tokenize arguments", logger.Err(err))
			return helpers.BadRequest(ctx, st.Message())
		}
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "tokenizerService.Tokenize failed",
			logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to tokenize")
	}

	token := hex.EncodeToString(tokenizeResp.TokenSuffix)
	if kind != nil && kind.ShortName != "" {
		token = kind.ShortName + "_" + token
	}

	if !pseudonymize {
		return ctx.JSON(http.StatusOK, &schemas.TokenizeResultSchema{Token: token})
	}

	mappingReq := &mapping.CreateMappingRequest{
		Token:         token,
		CipherText:    tokenizeResp.CipherText,
		DekWrapped:    tokenizeResp.DekWrapped,
		Deterministic: tokenizeResp.Deterministic,
		TokenTtl:      durationpb.New(time.Duration(tokenizeSchema.TokenTTL) * time.Second),
		AlgoName:      tokenizeResp.AlgoName,
	}
	if kind != nil {
		mappingReq.Kind = &mapping.Kind{Id: kind.Id}
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

	var kindID int32
	if kind != nil {
		kindID = kind.Id
	}
	if _, auditErr := t.mappingService.CreateAuditLog(reqCtx, &mapping.CreateAuditLogRequest{
		UserId: helpers.GetUserID(ctx),
		Action: "tokenize",
		Token:  token,
		KindId: kindID,
	}); auditErr != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to write audit log", logger.Err(auditErr))
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

	getMappingReq := &mapping.GetMappingByTokenRequest{Token: detokenizeSchema.Token}
	getMappingResp, err := t.mappingService.GetMappingByToken(reqCtx, getMappingReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "mapping not found",
					logger.Err(err))
				return helpers.NotFound(ctx, "token not found")
			case codes.InvalidArgument:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "invalid token",
					logger.Err(err))
				return helpers.BadRequest(ctx, "invalid arguments")
			case codes.DeadlineExceeded:
				logger.GetLoggerFromCtx(reqCtx).Info(reqCtx, "token expired",
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

	if kind := getMappingResp.MappingModel.Kind; kind != nil {
		if !helpers.HasRole(ctx, domain.RoleAdmin) && helpers.GetClearanceLevel(ctx) < int(kind.AccessLevel) {
			return helpers.Forbidden(ctx, "insufficient clearance level")
		}
	}

	detokenizeReq := &tokenizer.DetokenizeRequest{
		CipherText:    getMappingResp.MappingModel.CipherText,
		DekWrapped:    getMappingResp.MappingModel.DekWrapped,
		Deterministic: getMappingResp.MappingModel.Deterministic,
		AlgoName:      getMappingResp.MappingModel.AlgoName,
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

	var kindID int32
	if kind := getMappingResp.MappingModel.Kind; kind != nil {
		kindID = kind.Id
	}
	if _, auditErr := t.mappingService.CreateAuditLog(reqCtx, &mapping.CreateAuditLogRequest{
		UserId: helpers.GetUserID(ctx),
		Action: "detokenize",
		Token:  detokenizeSchema.Token,
		KindId: kindID,
	}); auditErr != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to write audit log", logger.Err(auditErr))
	}

	return ctx.JSON(http.StatusOK, &schemas.DetokenizeRespSchema{Plaintext: detokenizeResp.Plaintext})
}
