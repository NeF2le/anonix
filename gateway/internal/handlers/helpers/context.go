package helpers

import (
	"github.com/NeF2le/anonix/common/gen/auth_service"
	"github.com/NeF2le/anonix/gateway/internal/domain"
	"github.com/labstack/echo/v4"
)

func GetClearanceLevel(c echo.Context) int {
	level, ok := c.Get("clearanceLevel").(int)
	if !ok {
		return 1
	}
	return level
}

func GetUserID(c echo.Context) string {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return ""
	}
	return userID
}

func HasRole(c echo.Context, role domain.RoleName) bool {
	roles, ok := c.Get("roles").([]*auth_service.Role)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r.Name == string(role) {
			return true
		}
	}

	return false
}
