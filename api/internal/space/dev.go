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
	if !validateRequest(req) {
		return core.BadRequest(c, "missing fields in request body", nil)
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

func validateSend(m *core.Message) bool {
	if m == nil ||
		m.Data == nil ||
		m.Keyword == nil {
		return false
	}
	return true
}

func SendHandler(c echo.Context) error {
	req := new(core.Message)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid send format", err)
	}
	if !validateSend(req) {
		return core.BadRequest(c, "missing data or keyword", nil)
	}
	data, err := hex.DecodeString(*req.Data)
	if err != nil {
		return core.BadRequest(c, "invalid data format", err)
	}
	ciphertext, err := hex.DecodeString(*req.Keyword)
	if err != nil {
		return core.BadRequest(c, "invalid keyword format", err)
	}
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	from := claims.Username
	ok, space, err := db.FindSpace(c.Param("dev"))
	if !ok || err != nil {
		return core.ServerError(c, err)
	}
	serverKey := c.Get("ServerKey").(core.KeyContext)
	tag, err := db.MessageTag(ciphertext, serverKey.SKey, space)
	if err != nil {
		return core.ServerError(c, err)
	}
	ok, err = db.AddMessage(&db.Message{
		From:    from,
		To:      space.Owner,
		On:      space.Name,
		Tag:     tag,
		Data:    data,
		Keyword: ciphertext,
	})
	if !ok || err != nil {
		return core.ServerError(c, err)
	}
	return core.SendOK(c, nil)
}
