package auth

import (
	"encoding/hex"
	"errors"

	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"
	"github.com/labstack/echo/v4"
)

func validateRegister(u *core.User) bool {
	if u == nil ||
		u.Username == nil ||
		u.Password == nil ||
		u.PubKey == nil {
		return false
	}
	return true
}

func RegisterHandler(c echo.Context) error {
	req := new(core.User)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid request body", err)
	}
	if !validateRegister(req) {
		return core.BadRequest(c, "invalid request body", nil)
	}
	pubkey, err := hex.DecodeString(*req.PubKey)
	if err != nil {
		return core.BadRequest(c, "invalid public key", err)
	}
	u := &db.User{
		Username: *req.Username,
		Password: *req.Password,
		Pubkey:   pubkey,
	}
	ok, err := db.AddUser(u)
	if err != nil {
		return core.ServerError(c, err)
	}
	if !ok {
		return core.BadRequest(c, "invalid username", errors.New("username already exists"))
	}
	return core.SendOK(c, nil)
}
