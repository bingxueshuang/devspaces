package main

import (
	"log"
	"net/http"

	"github.com/bingxueshuang/devspaces/api/internal/auth"
	api "github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/bingxueshuang/devspaces/core"

	"github.com/labstack/echo/v4"
)

func main() {
	sk, pk, err := core.KeyGenServer()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
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
	e.GET("/", func(c echo.Context) error {
		return api.SendOK(c, "hello world")
	})
	authGroup := e.Group("/auth")
	auth.Setup(authGroup)
	e.GET("/user/:uname", auth.UserHandler)
	if err := e.Start(":5005"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
