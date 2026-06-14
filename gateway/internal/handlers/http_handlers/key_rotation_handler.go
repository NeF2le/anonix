package http_handlers

import (
	"context"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/common/gen/tokenizer"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/handlers/helpers"
	"github.com/NeF2le/anonix/gateway/internal/schemas"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type KeyRotationHandler struct {
	tokenizerService *services.TokenizerService
	mappingService   *services.MappingService
}

func NewKeyRotationHandler(
	tokenizerService *services.TokenizerService,
	mappingService *services.MappingService) *KeyRotationHandler {
	return &KeyRotationHandler{
		tokenizerService: tokenizerService,
		mappingService:   mappingService,
	}
}

// RotateMasterKey godoc
// @Summary Ротация мастер-ключа (KEK)
// @Description Создаёт новую версию мастер-ключа Vault Transit и перешифровывает обёртки DEK всех маппингов последней версией ключа. Данные не изменяются.
// @Tags Security
// @Produce json
// @Success 200 {object} schemas.KeyRotationResultSchema
// @Failure 500 "internal error"
// @Security ApiKeyAuth
// @Router /admin/keys/rotate-master [post]
func (k *KeyRotationHandler) RotateMasterKey(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	if _, err := k.tokenizerService.RotateMasterKey(reqCtx, &tokenizer.RotateMasterKeyRequest{}); err != nil {
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to rotate master key", logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to rotate master key")
	}

	listResp, err := k.mappingService.GetMappingList(reqCtx, &mapping.GetMappingListRequest{})
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to get mapping list", logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to get mapping list")
	}

	var updated, failed int32
	for _, mp := range listResp.GetMappingModels() {
		if rotateErr := k.rewrapMappingDek(reqCtx, mp); rotateErr != nil {
			logger.GetLoggerFromCtx(reqCtx).Debug(reqCtx, "failed to rewrap mapping dek",
				slog.String("id", mp.GetId()),
				logger.Err(rotateErr))
			failed++
			continue
		}
		updated++
	}

	if _, auditErr := k.mappingService.CreateAuditLog(reqCtx, &mapping.CreateAuditLogRequest{
		UserId: helpers.GetUserID(ctx),
		Action: "rotate_master_key",
	}); auditErr != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to write audit log", logger.Err(auditErr))
	}

	return ctx.JSON(http.StatusOK, &schemas.KeyRotationResultSchema{UpdatedCount: updated, FailedCount: failed})
}

func (k *KeyRotationHandler) rewrapMappingDek(ctx context.Context, mp *mapping.MappingModel) error {
	rewrapResp, err := k.tokenizerService.RewrapDEK(ctx, &tokenizer.RewrapDEKRequest{DekWrapped: mp.GetDekWrapped()})
	if err != nil {
		return err
	}

	_, err = k.mappingService.UpdateMappingDek(ctx, &mapping.UpdateMappingDekRequest{
		Id:         mp.GetId(),
		DekWrapped: rewrapResp.GetDekWrapped(),
	})
	return err
}

// RotateAllDeks godoc
// @Summary Ротация ключей шифрования данных (DEK)
// @Description Полностью перешифровывает данные всех маппингов новыми DEK. Значения токенов, видимые пользователям, не меняются.
// @Tags Security
// @Produce json
// @Success 200 {object} schemas.KeyRotationResultSchema
// @Failure 500 "internal error"
// @Security ApiKeyAuth
// @Router /admin/keys/rotate-deks [post]
func (k *KeyRotationHandler) RotateAllDeks(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	listResp, err := k.mappingService.GetMappingList(reqCtx, &mapping.GetMappingListRequest{})
	if err != nil {
		logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "failed to get mapping list", logger.Err(err))
		return helpers.InternalServerError(ctx, "failed to get mapping list")
	}

	var updated, failed int32
	for _, mp := range listResp.GetMappingModels() {
		if rotateErr := k.rotateMappingDek(reqCtx, mp); rotateErr != nil {
			logger.GetLoggerFromCtx(reqCtx).Debug(reqCtx, "failed to rotate mapping dek",
				slog.String("id", mp.GetId()),
				logger.Err(rotateErr))
			failed++
			continue
		}
		updated++
	}

	if _, auditErr := k.mappingService.CreateAuditLog(reqCtx, &mapping.CreateAuditLogRequest{
		UserId: helpers.GetUserID(ctx),
		Action: "rotate_deks",
	}); auditErr != nil {
		logger.GetLoggerFromCtx(reqCtx).Warn(reqCtx, "failed to write audit log", logger.Err(auditErr))
	}

	return ctx.JSON(http.StatusOK, &schemas.KeyRotationResultSchema{UpdatedCount: updated, FailedCount: failed})
}

func (k *KeyRotationHandler) rotateMappingDek(ctx context.Context, mp *mapping.MappingModel) error {
	rotateResp, err := k.tokenizerService.RotateDEK(ctx, &tokenizer.RotateDEKRequest{
		DekWrapped:    mp.GetDekWrapped(),
		CipherText:    mp.GetCipherText(),
		Deterministic: mp.GetDeterministic(),
		AlgoName:      mp.GetAlgoName(),
	})
	if err != nil {
		return err
	}

	_, err = k.mappingService.UpdateMappingCrypto(ctx, &mapping.UpdateMappingCryptoRequest{
		Id:         mp.GetId(),
		DekWrapped: rotateResp.GetDekWrapped(),
		CipherText: rotateResp.GetCipherText(),
		AlgoName:   rotateResp.GetAlgoName(),
	})
	return err
}
