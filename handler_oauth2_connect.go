package main

import (
	"flow-users/flags"
	"flow-users/jwt"
	"flow-users/oauth2"
	"flow-users/oauth2/github"
	"flow-users/oauth2/google"
	"flow-users/oauth2/twitter"
	"flow-users/user"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type OAuth2GitHubPost struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type OAuth2GooglePost struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type OAuth2TwitterPost struct {
	AccessToken          string `json:"access_token" validate:"required"`
	ExpireIn             int64  `json:"expire_in" validate:"required"`
	RefreshToken         string `json:"refresh_token" validate:"required"`
	RefreshTokenExpireIn int64  `json:"refresh_token_expire_in"`
}

func connectOAuth2(c echo.Context) (err error) {
	// Check token
	token := c.Get("user").(*jwtGo.Token)
	user_id, err := jwt.CheckToken(*flags.Get().JwtIssuer, token)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Get user
	u, notFound, err := user.Get(user_id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderGitHub.String():
		if *flags.Get().GithubClientId == "" || *flags.Get().GithubClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	case oauth2.ProviderGoogle.String():
		if *flags.Get().GoogleClientId == "" || *flags.Get().GoogleClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	case oauth2.ProviderTwitter.String():
		if *flags.Get().TwitterClientId == "" || *flags.Get().TwitterClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	default:
		// 404: Not found
		c.Logger().Debug("provider not found")
		return echo.ErrNotFound
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
		a, err := github.New(*flags.Get().GithubClientId, *flags.Get().GithubClientSecret)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Error(err)
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
			c.Logger().Error(err)
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
		a, err := google.New(*flags.Get().GoogleClientId, *flags.Get().GoogleClientSecret)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Write to DB
		_, err = google.Insert(
			google.OAuth2{
				AccessToken: p.AccessToken,
				OwnerId:     o.Id,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Error(err)
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
		a, err := twitter.New(*flags.Get().TwitterClientId, *flags.Get().TwitterClientSecret)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Error(err)
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
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
	}

	// Generate token
	t, err := jwt.GenerateToken(user.UserWithoutPassword{Id: u.Id, Name: u.Name, Email: u.Email}, *flags.Get().JwtIssuer, *flags.Get().JwtSecret)
	if err != nil {
		c.Logger().Error(err)
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
