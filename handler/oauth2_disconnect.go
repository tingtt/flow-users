package handler

import (
	"flow-users/flags"
	"flow-users/jwt"
	"flow-users/oauth2"
	"flow-users/oauth2/github"
	"flow-users/oauth2/google"
	"flow-users/oauth2/twitter"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func DisconnectOAuth2(c echo.Context) (err error) {
	// Check token
	token := c.Get("user").(*jwtGo.Token)
	user_id, err := jwt.CheckToken(*flags.Get().JwtIssuer, token)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
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
		c.Logger().Debugf("provider '%s' not found", provider)
		return echo.ErrNotFound
	}

	switch provider {
	case "github":
		// Write to DB
		notFound, err := github.Delete(user_id)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			c.Logger().Debug("GitHub OAuth2 connection not found")
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "GitHub OAuth2 connection not found"}, "	")
		}

	case "google":
		// Write to DB
		notFound, err := google.Delete(user_id)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			c.Logger().Debug("Google OAuth2 connection not found")
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Google OAuth2 connection not found"}, "	")
		}

	case "twitter":
		// Write to DB
		notFound, err := twitter.Delete(user_id)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			c.Logger().Debug("Twitter OAuth2 connection not found")
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Twitter OAuth2 connection not found"}, "	")
		}
	}

	// 204: No content
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
