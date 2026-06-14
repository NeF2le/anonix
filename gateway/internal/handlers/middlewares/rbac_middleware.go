package middlewares

import (
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/gateway/internal/domain"
	"github.com/labstack/echo/v4"
)

type RBACMiddleware struct {
}

func NewRBACMiddleware() *RBACMiddleware {
	return &RBACMiddleware{}
}

func (a *RBACMiddleware) CheckRole(required ...domain.RoleName) echo.MiddlewareFunc {
	requiredSet := make(map[string]struct{}, len(required))
	for _, r := range required {
		requiredSet[string(r)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := c.Get("roles")
			if val == nil {
				return echo.ErrForbidden
			}

			userRoles, ok := val.([]*auth_service.Role)
			if !ok {
				return echo.ErrInternalServerError
			}

			for _, r := range userRoles {
				if _, ok = requiredSet[r.Name]; ok {
					return next(c)
				}
			}

			return echo.ErrForbidden
		}
	}
}
