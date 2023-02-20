package space

import (
	"encoding/hex"

	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func validateSpace(s *core.DevSpace) bool {
	if s == nil ||
		s.Name == nil ||
		s.Pubkey == nil {
		return false
	}
	return true
}

func CreateDev(c echo.Context) error {
	req := new(core.DevSpace)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid request body", err)
	}
	pubkey, err := hex.DecodeString(*req.Pubkey)
	if err != nil {
		return core.BadRequest(c, "invalid public key", err)
	}
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	owner := claims.Username
	ok, err := db.AddSpace(&db.Space{
		Name:   *req.Name,
		Owner:  owner,
		Pubkey: pubkey,
		Tags:   nil,
	})
	if !ok || err != nil {
		return core.ServerError(c, err)
	}
	return core.SendOK(c, nil)
}

func ListDev(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*core.TokenClaims)
	owner := claims.Username
	spaces, err := db.ListSpaces(owner)
	if err != nil {
		return core.ServerError(c, err)
	}
	res := make([]map[string]any, len(spaces))
	for _, s := range spaces {
		res = append(res, map[string]any{
			"name":   s.Name,
			"pubkey": hex.EncodeToString(s.Pubkey),
		})
	}
	return core.SendOK(c, res)
}
