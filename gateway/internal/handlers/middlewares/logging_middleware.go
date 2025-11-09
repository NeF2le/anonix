package middlewares

import (
	"github.com/NeF2le/anonix/common/logger"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := logger.New(c.Request().Context())

		c.SetRequest(c.Request().WithContext(ctx))

		logger.GetLoggerFromCtx(ctx).Info(ctx, "HTTP request",
			slog.String("method", c.Request().Method),
			slog.String("path", c.Request().URL.Path))

		err := next(c)

		logger.GetLoggerFromCtx(ctx).Info(ctx, "Completed HTTP request",
			slog.String("method", c.Request().Method),
			slog.String("path", c.Request().URL.Path))

		return err
	}
}
