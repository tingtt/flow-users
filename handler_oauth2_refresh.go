package main

import (
	"flow-users/jwt"
	"flow-users/oauth2"
	"flow-users/oauth2/twitter"
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
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderTwitter.String():
		if *twitterClientId == "" || *twitterClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

		a, err := twitter.New(*twitterClientId, *twitterClientSecret)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// Read DB row
		owner, notFound, err := twitter.Get(user_id)
		if err != nil {
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			return echo.ErrNotFound
		}

		// Refresh token
		newOwner, err := a.RefreshToken(owner.RefreshToken)
		if err != nil {
			c.Logger().Error(err)
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
			c.Logger().Error(err)
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
