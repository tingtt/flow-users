package main

import (
	"flow-user/jwt"
	"flow-user/oauth2"
	"flow-user/oauth2/github"
	"flow-user/oauth2/google"
	"flow-user/oauth2/twitter"
	"flow-user/user"
	"net/http"

	"github.com/labstack/echo"
)

type UserPostOverOAuth2 struct {
	AccessToken          string `json:"access_token" validate:"require"`
	ExpireIn             int64  `json:"expire_in" validate:"require"`
	RefreshToken         string `json:"refresh_token" validate:"require"`
	RefreshTokenExpireIn int64  `json:"refresh_token_expire_in"`
	Password             string `json:"password" validate:"require"`
}

func postOverOAuth2(c echo.Context) (err error) {
	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderGitHub.String():
		if *githubClientId == "" || *githubClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	case oauth2.ProviderGoogle.String():
		if *googleClientId == "" || *googleClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	case oauth2.ProviderTwitter.String():
		if *twitterClientId == "" || *twitterClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	default:
		// 404: Not found
		c.Logger().Debug(err)
		return echo.ErrNotFound
	}

	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" {
		// 415: Invalid `Content-Type`
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	p := new(UserPostOverOAuth2)
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

	// Get owner info
	var u user.User
	var name string
	var email string

	switch provider {
	case "github":
		// Get owner info
		a, err := github.New(*githubClientId, *githubClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		name = o.Name

		e, err := a.GetOwnerPrimaryEmail(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		email = e.Email

		// Write to DB
		u, invalidEmail, usedEmail, err := user.Insert(user.UserPost{Name: name, Email: email, Password: p.Password})
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
			c.Logger().Debug("email already used")
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
		}

		// Write to DB
		_, err = github.Insert(
			github.OAuth2{
				AccessToken:          p.AccessToken,
				ExpireIn:             p.ExpireIn,
				RefreshToken:         p.RefreshToken,
				RefreshTokenExpireIn: p.RefreshTokenExpireIn,
				OwnerId:              o.Id,
			},
			u.Id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case "google":
		// Get owner info
		a, err := google.New(*googleClientId, *googleClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		name = o.Name
		email = o.Email

		// Write to DB
		u, invalidEmail, usedEmail, err := user.Insert(user.UserPost{Name: name, Email: email, Password: p.Password})
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
			c.Logger().Debug("email already used")
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
		}

		// Write to DB
		_, err = google.Insert(
			google.OAuth2{
				AccessToken:  p.AccessToken,
				ExpireIn:     p.ExpireIn,
				RefreshToken: p.RefreshToken,
				OwnerId:      o.Id,
			},
			u.Id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case "twitter":
		// Get owner info
		a, err := twitter.New(*googleClientId, *googleClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		name = o.UserName

		email, err = a.GetOwnerEmail(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Write to DB
		u, invalidEmail, usedEmail, err := user.Insert(user.UserPost{Name: name, Email: email, Password: p.Password})
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
			c.Logger().Debug("email already used")
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
		}

		// Write to DB
		_, err = twitter.Insert(
			twitter.OAuth2{
				AccessToken:          p.AccessToken,
				ExpireIn:             p.ExpireIn,
				RefreshToken:         p.RefreshToken,
				RefreshTokenExpireIn: p.RefreshTokenExpireIn,
				OwnerId:              o.Id,
			},
			u.Id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
	}

	// Generate token
	t, err := jwt.GenerateToken(user.UserPostResponse{Id: u.Id, Name: name, Email: email}, *jwtIssuer, *jwtSecret)
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
	return c.JSONPretty(http.StatusOK, user.UserPostResponse{Id: u.Id, Name: name, Email: email}, "	")
}
