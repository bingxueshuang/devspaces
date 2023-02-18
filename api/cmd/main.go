package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world")
	})
	if err := e.Start(":5005"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
