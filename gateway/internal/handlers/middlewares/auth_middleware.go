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

		sub, finalAccess, roles, clearanceLevel, err := a.ensureValidAccessToken(c, accessToken)
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
		c.Set("roles", roles)
		c.Set("clearanceLevel", clearanceLevel)

		return next(c)
	}
}

func (a *AuthMiddleware) setAuthHeader(c echo.Context, token string) {
	c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (a *AuthMiddleware) ensureValidAccessToken(c echo.Context, accessToken string) (string, string, []*auth_service.Role, int, error) {
	if accessToken == "" {
		newAccess, err := a.refreshAndSetCookies(c)
		if err != nil {
			return "", "", nil, 0, errs.ErrUnauthorized
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	sub, isRefresh, expTime, roles, clearanceLevel, err := a.parseJwt(accessToken, a.jwtSecret)
	if err != nil {
		newAccess, refErr := a.refreshAndSetCookies(c)
		if refErr != nil {
			return "", "", nil, 0, errs.ErrUnauthorized
		}
		sub, isRefresh, expTime, roles, clearanceLevel, err = a.parseJwt(newAccess, a.jwtSecret)
		if err != nil {
			return "", "", nil, 0, errs.ErrUnauthorized
		}
		if isRefresh {
			return "", "", nil, 0, errs.ErrInvalidAccessToken
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	if isRefresh {
		return "", "", nil, 0, errs.ErrInvalidAccessToken
	}

	if time.Now().After(expTime) {
		newAccess, refErr := a.refreshAndSetCookies(c)
		if refErr != nil {
			return "", "", nil, 0, errs.ErrUnauthorized
		}
		sub, isRefresh, expTime, roles, clearanceLevel, err = a.parseJwt(newAccess, a.jwtSecret)
		if err != nil {
			return "", "", nil, 0, errs.ErrUnauthorized
		}
		if isRefresh {
			return "", "", nil, 0, errs.ErrInvalidAccessToken
		}
		accessToken = newAccess
		a.setAuthHeader(c, accessToken)
	}

	_ = expTime
	return sub, accessToken, roles, clearanceLevel, nil
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

func (a *AuthMiddleware) parseJwt(rawToken, jwtSecret string) (string, bool, time.Time, []*auth_service.Role, int, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", false, time.Time{}, nil, 0, errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, time.Time{}, nil, 0, errs.ErrInvalidToken
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return "", false, time.Time{}, nil, 0, errs.ErrInvalidToken
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", false, time.Time{}, nil, 0, err
	}

	isRefresh, _ := claims["is_refresh"].(bool)

	var roles []*auth_service.Role
	if raw, ok := claims["roles"].([]interface{}); ok {
		for _, r := range raw {
			if name, ok := r.(string); ok {
				roles = append(roles, &auth_service.Role{Name: name})
			}
		}
	}

	clearanceLevel := 1
	if raw, ok := claims["clearance_level"].(float64); ok {
		clearanceLevel = int(raw)
	}

	return sub, isRefresh, time.Unix(exp.Unix(), 0), roles, clearanceLevel, nil
}
