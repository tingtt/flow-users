package main

import (
	"errors"
	"flow-user/jwt"
	"flow-user/oauth2"
	"flow-user/oauth2/github"
	"flow-user/oauth2/google"
	"flow-user/oauth2/twitter"
	"flow-user/user"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type OAuth2GitHubPost struct {
	AccessToken string `json:"access_token" validate:"require"`
}

type OAuth2GooglePost struct {
	AccessToken  string `json:"access_token" validate:"require"`
	ExpireIn     int64  `json:"expire_in" validate:"require"`
	RefreshToken string `json:"refresh_token" validate:"require"`
}

type OAuth2TwitterPost struct {
	AccessToken          string `json:"access_token" validate:"require"`
	ExpireIn             int64  `json:"expire_in" validate:"require"`
	RefreshToken         string `json:"refresh_token" validate:"require"`
	RefreshTokenExpireIn int64  `json:"refresh_token_expire_in"`
}

func connectOAuth2(c echo.Context) (err error) {
	// Check token
	token := c.Get("user").(*jwtGo.Token)
	user_id, err := jwt.CheckToken(*jwtIssuer, token)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Get user
	u, notFound, err := user.Get(user_id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

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
		c.Logger().Debug(errors.New("provider not found"))
		return echo.ErrNotFound
	}

	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" &&
		c.Request().Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// 415: Invalid `Content-Type`
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()}, "	")
	}

	switch provider {
	case "github":
		// Bind request body
		p := new(OAuth2GitHubPost)
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

		// Write to DB
		_, err = github.Insert(
			github.OAuth2{
				AccessToken: p.AccessToken,
				OwnerId:     o.Id,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case "google":
		// Bind request body
		p := new(OAuth2GooglePost)
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

		// Write to DB
		_, err = google.Insert(
			google.OAuth2{
				AccessToken:  p.AccessToken,
				ExpireIn:     p.ExpireIn,
				RefreshToken: p.RefreshToken,
				OwnerId:      o.Id,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case "twitter":
		// Bind request body
		p := new(OAuth2TwitterPost)
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
		a, err := twitter.New(*twitterClientId, *twitterClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
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
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
	}

	// Generate token
	t, err := jwt.GenerateToken(user.UserPostResponse{Id: u.Id, Name: u.Name, Email: u.Email}, *jwtIssuer, *jwtSecret)
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
	return c.JSONPretty(http.StatusOK, map[string]string{"message": "Success"}, "	")
}
