package space

import "github.com/labstack/echo/v4"

func Setup(g *echo.Group) {
	g.POST("/:dev/request", RequestHandler)
}
