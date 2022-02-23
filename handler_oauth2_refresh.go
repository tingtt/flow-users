package main

import (
	"flow-user/jwt"
	"flow-user/oauth2"
	"flow-user/oauth2/github"
	"flow-user/oauth2/google"
	"flow-user/oauth2/twitter"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func refreshOAuth2Token(c echo.Context) (err error) {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	user_id, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderGitHub.String():
		if *githubClientId == "" || *githubClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

		a, err := github.New(*githubClientId, *githubClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Read DB row
		owner, notFound, err := github.Get(user_id)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			return echo.ErrNotFound
		}

		// Refresh token
		newOwner, err := a.RefreshToken(owner.RefreshToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Update DB row
		_, err = github.Insert(
			github.OAuth2{
				AccessToken:          newOwner.AccessToken,
				ExpireIn:             newOwner.ExpireIn,
				RefreshToken:         newOwner.RefreshToken,
				RefreshTokenExpireIn: newOwner.RefreshTokenExpireIn,
				OwnerId:              owner.OwnerId,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case oauth2.ProviderGoogle.String():
		if *googleClientId == "" || *googleClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

		a, err := google.New(*googleClientId, *googleClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Read DB row
		owner, notFound, err := google.Get(user_id)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			return echo.ErrNotFound
		}

		// Refresh token
		newOwner, err := a.RefreshToken(owner.RefreshToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Update DB row
		_, err = google.Insert(
			google.OAuth2{
				AccessToken:  newOwner.AccessToken,
				ExpireIn:     newOwner.ExpireIn,
				RefreshToken: newOwner.RefreshToken,
				OwnerId:      owner.OwnerId,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	case oauth2.ProviderTwitter.String():
		if *twitterClientId == "" || *twitterClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

		a, err := twitter.New(*twitterClientId, *twitterClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Read DB row
		owner, notFound, err := twitter.Get(user_id)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			return echo.ErrNotFound
		}

		// Refresh token
		newOwner, err := a.RefreshToken(owner.RefreshToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Update DB row
		_, err = twitter.Insert(
			twitter.OAuth2{
				AccessToken:          newOwner.AccessToken,
				ExpireIn:             newOwner.ExpireIn,
				RefreshToken:         newOwner.RefreshToken,
				RefreshTokenExpireIn: newOwner.RefreshTokenExpireIn,
				OwnerId:              owner.OwnerId,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

	default:
		// 404: Not found
		c.Logger().Debug(err)
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, map[string]string{"message": "Success"}, "	")
}
