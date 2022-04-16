package main

import (
	"errors"
	"flow-users/jwt"
	"flow-users/user"
	"net/http"

	"github.com/labstack/echo"
)

func post(c echo.Context) (err error) {
	// Bind request body
	p := new(user.PostBody)
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

	// Write to DB
	u, invalidEmail, usedEmail, err := user.Post(*p)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if invalidEmail {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}
	if usedEmail {
		// 409: Conflict
		c.Logger().Debug(errors.New("email already used"))
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
	}

	// Generate token
	t, err := jwt.GenerateToken(p.PostResponse(u.Id), *jwtIssuer, *jwtSecret)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		HttpOnly: true,
	})

	// 200: Success
	return c.JSONPretty(http.StatusOK, p.PostResponse(u.Id), "	")
}
