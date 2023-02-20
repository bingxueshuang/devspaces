package space

import (
	"encoding/hex"

	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func DashboardHandler(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	uname := claims.Username
	requests, err := db.RequestsTo(uname)
	if err != nil {
		return core.ServerError(c, err)
	}
	res := make([]map[string]any, len(requests))
	for _, r := range requests {
		res = append(res, map[string]any{
			"from":   r.From,
			"on":     r.On,
			"to":     r.To,
			"secret": hex.EncodeToString(r.Secret),
		})
	}
	return core.SendOK(c, res)
}
