package twitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Authentication struct {
	AccessToken          string `json:"access_token"`
	ExpireIn             int64  `json:"expire_in"`
	RefreshToken         string `json:"refresh_token"`
	RefreshTokenExpireIn int64  `json:"refresh_token_expire_in"`
}

type refreshPostBody struct {
	AccessToken  string `json:"access_token"`
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (g *Application) RefreshToken(refreshToken string) (owner Authentication, err error) {
	// Create request body
	j, err := json.Marshal(refreshPostBody{
		AccessToken:  refreshToken,
		GrantType:    "refresh_token",
		ClientId:     g.ClientId,
		ClientSecret: g.ClientSecret,
	})
	if err != nil {
		return Authentication{}, errors.New("failed to create body using refresh twitter access_token")
	}

	// POST twitter api
	res, err := http.Post("https://api.twitter.com/2/oauth2/token", "application/json", bytes.NewBuffer(j))
	if err != nil {
		return Authentication{}, errors.New("failed to refresh token")
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return Authentication{}, errors.New("failed to refresh token")
	}

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Authentication{}, errors.New("failed to read body of response from twitter api")
	}

	// Unmarshal response body
	err = json.Unmarshal(bodyBytes, &owner)
	if err != nil {
		return Authentication{}, errors.New("failed to read body of response from twitter api")
	}

	return owner, nil
}
