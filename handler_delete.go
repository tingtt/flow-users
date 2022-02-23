package main

import (
	"errors"
	"flow-user/jwt"
	"flow-user/user"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func delete(c echo.Context) (err error) {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	id, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Delete DB row
	notFound, err := user.Delete(id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug(errors.New("user not found"))
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
