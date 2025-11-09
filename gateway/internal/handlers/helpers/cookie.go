package helpers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SetAccessTokenCookie(c echo.Context, accessToken string, maxAge int) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return
}

func SetRefreshTokenCookie(c echo.Context, refreshToken string, maxAge int) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return
}
