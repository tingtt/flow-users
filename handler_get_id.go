package main

import (
	"flow-users/flags"
	"flow-users/jwt"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func getId(c echo.Context) (err error) {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	id, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, map[string]uint64{"id": id}, "	")
}
