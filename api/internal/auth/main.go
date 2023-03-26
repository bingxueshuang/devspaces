package auth

import (
	"github.com/labstack/echo/v4"
)

func Setup(g *echo.Group) {
	g.POST("/register", RegisterHandler)
	g.POST("/login", LoginHandler)
	//g.GET("/debug", func(c echo.Context) error {
	//	return core.SendOK(c, db.ListUsers())
	//})
}
