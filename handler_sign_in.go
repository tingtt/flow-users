package main

import (
	"flow-users/jwt"
	"flow-users/user"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func signIn(c echo.Context) (err error) {
	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" &&
		c.Request().Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// 415: Invalid `Content-Type`
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	p := new(user.UserSignInPost)
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

	// Get user by email and compare password
	u, notFound, err := user.GetByEmail(p.Email)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// Incorrect email
		// 404: Not found
		return echo.ErrNotFound
	}
	verify, err := u.Verify(p.Password)
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if !verify {
		// Incorrect password
		// 403: Forbidden
		c.Logger().Debug("failed to sign in")
		return echo.ErrForbidden
	}

	// Generate token
	t, err := jwt.GenerateToken(
		user.UserPostResponse{Id: u.Id, Name: u.Name, Email: u.Email},
		*jwtIssuer,
		*jwtSecret,
	)
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
	return c.JSONPretty(
		http.StatusOK,
		map[string]string{"token": t},
		"	",
	)
}
