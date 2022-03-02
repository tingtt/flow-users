package main

import (
	"flow-users/jwt"
	"flow-users/user"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func patch(c echo.Context) (err error) {
	// Check token
	t := c.Get("user").(*jwtGo.Token)
	user_id, err := jwt.CheckToken(*jwtIssuer, t)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" &&
		c.Request().Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Bind request body
	p := new(user.UserPatch)
	if err = c.Bind(p); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(p); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Update DB row
	u, invalidEmail, usedEmail, notFound, err := user.Update(user_id, *p)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if invalidEmail {
		// 422: Unprocessable entity
		c.Logger().Debug("invalid email")
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": "invalid email"}, "	")
	}
	if usedEmail {
		// 409: Conflict
		c.Logger().Debug("email already used")
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("user not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, u, "	")
}
