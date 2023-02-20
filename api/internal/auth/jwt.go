package auth

import (
	"time"

	"github.com/bingxueshuang/devspaces/api/internal/core"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/labstack/echo-jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func getToken(username string) (string, error) {
	claims := core.TokenClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// create a token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it
	return token.SignedString(core.TokenSecret)
}

var Config = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(core.TokenClaims)
	},
	SigningKey: core.TokenSecret,
}
