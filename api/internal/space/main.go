package space

import "github.com/labstack/echo/v4"

func Setup(g *echo.Group) {
	g.POST("/:dev/request", RequestHandler)
	g.POST("/:dev/send", SendHandler)
	g.POST("/", CreateDev)
	g.GET("/", ListDev)
	g.POST("/:dev", CreateTag)
	g.GET("/:dev", ListTags)
	g.GET("/:dev/pubkey", PubkeyHandler)
	g.GET("/:dev/:tag", ListMessages)
}
