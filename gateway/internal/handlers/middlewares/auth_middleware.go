package middlewares

import (
	"errors"
	"fmt"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/NeF2le/anonix/gateway/internal/handlers/helpers"
	"github.com/NeF2le/anonix/gateway/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"strings"
	"time"
)

type AuthMiddleware struct {
	jwtSecret       string
	authService     *services.AuthService
	accessTokenTTL  int
	refreshTokenTTL int
}

func NewAuthMiddleware(
	jwtSecret string,
	authService *services.AuthService,
	accessTokenTTL,
	refreshTokenTTL int) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:       jwtSecret,
		authService:     authService,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (a *AuthMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqCtx := c.Request().Context()
		accessToken := a.getAccessTokenFromRequest(c)
		logger.GetLoggerFromCtx(reqCtx).Debug(reqCtx, "got access token",
			slog.String("access_token", accessToken))

		sub, finalAccess, err := a.ensureValidAccessToken(c, accessToken)
		if err != nil {
			switch {
			case errors.Is(err, errs.ErrInvalidAccessToken):
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "invalid access token",
					logger.Err(err))
				return helpers.BadRequest(c, "invalid access token")
			default:
				logger.GetLoggerFromCtx(reqCtx).Error(reqCtx, "error to get access token",
					logger.Err(err))
				return helpers.Unauthorized(c)
			}
		}

		if finalAccess != "" {
			a.setAuthHeader(c, finalAccess)
		}
		c.Set("userID", sub)
		return next(c)
	}
}

func (a *AuthMiddleware) setAuthHeader(c echo.Context, token string) {
	c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (a *AuthMiddleware) ensureValidAccessToken(c echo.Context, accessToken string) (string, string, error) {
	if accessToken == "" {
		newAccess, err := a.refreshAndSetCookies(c)
		if err != nil {
			return "", "", errs.ErrUnauthorized
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	// Parsing access token
	sub, isRefresh, expTime, err := a.parseJwt(accessToken, a.jwtSecret)
	if err != nil {
		// If parsing failed, trying to refresh
		newAccess, refErr := a.refreshAndSetCookies(c)
		if refErr != nil {
			return "", "", errs.ErrUnauthorized
		}
		sub, isRefresh, expTime, err = a.parseJwt(newAccess, a.jwtSecret)
		if err != nil {
			return "", "", errs.ErrUnauthorized
		}
		if isRefresh {
			return "", "", errs.ErrInvalidAccessToken
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	if isRefresh {
		return "", "", errs.ErrInvalidAccessToken
	}

	if time.Now().After(expTime) {
		newAccess, refErr := a.refreshAndSetCookies(c)
		if refErr != nil {
			return "", "", errs.ErrUnauthorized
		}
		sub, isRefresh, expTime, err = a.parseJwt(newAccess, a.jwtSecret)
		if err != nil {
			return "", "", errs.ErrUnauthorized
		}
		if isRefresh {
			return "", "", errs.ErrInvalidAccessToken
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	return sub, accessToken, nil
}

func (a *AuthMiddleware) getAccessTokenFromRequest(c echo.Context) string {
	var accessToken string
	auth := c.Request().Header.Get("Authorization")
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		accessToken = strings.TrimPrefix(auth, "Bearer ")
	} else {
		if ck, err := c.Cookie("access_token"); err == nil && ck.Value != "" {
			accessToken = ck.Value
		}
	}
	return accessToken
}

func (a *AuthMiddleware) refreshAndSetCookies(c echo.Context) (string, error) {
	rck, err := c.Cookie("refresh_token")
	if err != nil || rck.Value == "" {
		return "", err
	}

	req := &auth_service.RefreshRequest{RefreshToken: rck.Value}
	resp, err := a.authService.Refresh(c.Request().Context(), req)
	if err != nil {
		return "", err
	}

	helpers.SetAccessTokenCookie(c, resp.AccessToken, a.accessTokenTTL)
	helpers.SetRefreshTokenCookie(c, resp.RefreshToken, a.refreshTokenTTL)

	return resp.AccessToken, nil
}

func (a *AuthMiddleware) parseJwt(rawToken, jwtSecret string) (string, bool, time.Time, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", false, time.Time{}, errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, time.Time{}, errs.ErrInvalidToken
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return "", false, time.Time{}, errs.ErrInvalidToken
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", false, time.Time{}, err
	}

	expTime := time.Unix(exp.Unix(), 0)

	isRefresh, _ := claims["is_refresh"].(bool)

	return sub, isRefresh, expTime, nil
}
