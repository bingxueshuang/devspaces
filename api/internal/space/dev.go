package space

import (
	"encoding/hex"
	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func validateRequest(r *core.Request) bool {
	if r == nil ||
		r.To == nil ||
		r.Secret == nil {
		return false
	}
	return true
}

func RequestHandler(c echo.Context) error {
	req := new(core.Request)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid collaboration request", err)
	}
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	from := claims.Username
	secret, err := hex.DecodeString(*req.Secret)
	if err != nil {
		return core.BadRequest(c, "invalid request secret", err)
	}
	devspace := c.Param("dev")
	ok, err := db.AddRequest(&db.Request{
		From:   from,
		On:     devspace,
		To:     *req.To,
		Secret: secret,
	})
	if !ok || err != nil {
		return core.ServerError(c, err)
	}
	return core.SendOK(c, nil)
}
