package main

import (
	"encoding/hex"
	"log"
	"net/http"

	"github.com/bingxueshuang/devspaces/api/internal/space"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/bingxueshuang/devspaces/api/internal/auth"
	api "github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/core"

	"github.com/labstack/echo/v4"
)

func PubkeyHandler(c echo.Context) error {
	serverKey := c.Get("ServerKey").(api.KeyContext)
	pk := serverKey.PKey.Bytes()
	return api.SendOK(c, map[string]any{
		"pubkey": hex.EncodeToString(pk),
	})
}

func main() {
	sk, pk, err := core.KeyGenServer()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.HideBanner = true
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			kc := api.KeyContext{
				Context: c,
				SKey:    sk,
				PKey:    pk,
			}
			c.Set("ServerKey", kc)
			return next(c)
		}
	})
	authGroup := e.Group("/auth")
	auth.Setup(authGroup)
	e.GET("/user/:uname", auth.UserHandler)
	ptdGroup := e.Group("/space", echojwt.WithConfig(auth.Config))
	space.Setup(ptdGroup)
	e.GET("/dashboard", space.DashboardHandler, echojwt.WithConfig(auth.Config))
	e.GET("/pubkey", PubkeyHandler)
	e.GET("/", func(c echo.Context) error {
		return api.SendOK(c, "hello world")
	})

	if err := e.Start(":5005"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
