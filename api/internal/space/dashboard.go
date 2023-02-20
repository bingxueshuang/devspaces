package space

import (
	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func DashboardHandler(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	uname := claims.Username
	return core.SendOK(c, "welcome "+uname)
}
