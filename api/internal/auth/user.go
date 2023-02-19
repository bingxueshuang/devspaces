package auth

import (
	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"
	"github.com/labstack/echo/v4"
)

func UserHandler(c echo.Context) error {
	username := c.Param("uname")
	ok, user, err := db.GetUser(username)
	if err != nil {
		return core.ServerError(c, err)
	}
	if !ok {
		return core.NotFound(c, "username do not exist")
	}
	return core.SendOK(c, map[string]any{
		"username": user.Username,
		"pubkey":   user.Pubkey,
	})
}
