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
	if !validateSpace(req) {
		return core.BadRequest(c, "missing fields in request body", nil)
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
	res := make([]map[string]any, 0, len(spaces))
	for _, s := range spaces {
		res = append(res, map[string]any{
			"name":   s.Name,
			"pubkey": hex.EncodeToString(s.Pubkey),
		})
	}
	return core.SendOK(c, res)
}

func validateTag(t *core.Tag) bool {
	if t == nil ||
		t.Name == nil ||
		t.Trapdoor == nil {
		return false
	}
	return true
}

func CreateTag(c echo.Context) error {
	req := new(core.Tag)
	if err := c.Bind(req); err != nil {
		return core.BadRequest(c, "invalid request body", err)
	}
	if !validateTag(req) {
		return core.BadRequest(c, "missing fields in request body", nil)
	}
	trapdoor, err := hex.DecodeString(*req.Trapdoor)
	if err != nil {
		return core.BadRequest(c, "invalid trapdoor", err)
	}
	space := c.Param("dev")
	ok, err := db.AddTag(space, &db.Tag{
		Name:     *req.Name,
		Trapdoor: trapdoor,
	})
	if !ok || err != nil {
		return core.ServerError(c, err)
	}
	return core.SendOK(c, nil)
}

func ListTags(c echo.Context) error {
	space := c.Param("dev")
	tags, err := db.ListTags(space)
	if err != nil {
		return core.ServerError(c, err)
	}
	res := make([]map[string]any, 0, len(tags))
	for _, t := range tags {
		res = append(res, map[string]any{
			"name":     t.Name,
			"trapdoor": hex.EncodeToString(t.Trapdoor),
		})
	}
	return core.SendOK(c, res)
}

func ListMessages(c echo.Context) error {
	space := c.Param("dev")
	tag := c.Param("tag")
	mlist, err := db.ListMessages(tag, space)
	if err != nil {
		return core.ServerError(c, err)
	}
	msgs := make([]map[string]any, 0, len(mlist))
	for _, m := range mlist {
		msgs = append(msgs, map[string]any{
			"from":    m.From,
			"data":    hex.EncodeToString(m.Data),
			"keyword": hex.EncodeToString(m.Keyword),
		})
	}
	return core.SendOK(c, msgs)
}
