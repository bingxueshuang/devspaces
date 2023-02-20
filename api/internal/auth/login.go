package auth

import (
	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"

	"github.com/labstack/echo/v4"
)

func validateLogin(u *core.User) bool {
	if u == nil ||
		u.Username == nil ||
		u.Password == nil {
		return false
	}
	return true
}

func LoginHandler(c echo.Context) error {
	req := new(core.User)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid request body", err)
	}
	if !validateLogin(req) {
		return core.BadRequest(c, "invalid request body", nil)
	}
	u := &db.User{
		Username: *req.Username,
		Password: *req.Password,
	}
	ok, err := db.MatchUser(u)
	if err != nil {
		return core.ServerError(c, err)
	}
	if !ok {
		return core.BadRequest(c, "invalid username or password", nil)
	}
	token, err := getToken(*req.Username)
	if err != nil {
		return core.ServerError(c, err)
	}
	return core.SendOK(c, map[string]string{
		"token": token,
	})
}
