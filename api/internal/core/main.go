package core

import (
	"net/http"

	"github.com/bingxueshuang/devspaces/core"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type KeyContext struct {
	echo.Context
	SKey *core.SKey
	PKey *core.PKeyServer
}

type Request struct {
	From   *string `json:"from"`
	On     *string `json:"on"`
	To     *string `json:"to"`
	Secret *string `json:"secret"`
}

type DevSpace struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Owner  string   `json:"owner"`
	Pubkey string   `json:"pubkey"`
}

type Response struct {
	Ok    bool `json:"ok"`
	Data  any  `json:"data"`
	Error any  `json:"error"`
}

type User struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	PubKey   *string `json:"pubkey"`
}

type TokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var TokenSecret = []byte("secret message")

func SendOK(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, Response{
		Ok:    true,
		Data:  data,
		Error: nil,
	})
}

func BadRequest(c echo.Context, msg string, err error) error {
	var e any
	if err != nil {
		e = err.Error()
	}
	return c.JSON(http.StatusBadRequest, Response{
		Ok:    false,
		Data:  e,
		Error: msg,
	})
}

func NotFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, Response{
		Ok:    false,
		Data:  nil,
		Error: msg,
	})
}

func ServerError(c echo.Context, err error) error {
	var msg any
	if err != nil {
		msg = err.Error()
	}
	return c.JSON(http.StatusInternalServerError, Response{
		Ok:    false,
		Data:  msg,
		Error: "sorry, could not process your request",
	})
}
