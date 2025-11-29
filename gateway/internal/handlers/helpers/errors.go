package helpers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func InternalServerError(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err})
}

func BadRequest(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err})
}

func NotFound(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusNotFound, map[string]string{"error": err})
}

func Conflict(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusConflict, map[string]string{"error": err})
}

func Unauthorized(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "please log in first"})
}

func InvalidCredentials(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
}
