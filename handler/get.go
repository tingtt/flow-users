package handler

import (
	"flow-users/flags"
	"flow-users/jwt"
	"flow-users/user"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Get(c echo.Context) (err error) {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	id, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	u2, notFound, err := user.GetWithoutPassword(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("user not found")
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, u2, "	")
}
