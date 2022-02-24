package google

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Owner struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	PictureUrl string `json:"picture"`
}

func (g *Application) GetOwner(token string) (o Owner, err error) {
	// GET google api
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return Owner{}, errors.New("failed to get owner informations")
	}
	param := req.URL.Query()
	param.Add("oauth_token", token)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return Owner{}, err
	}
	defer res.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Owner{}, errors.New("failed to read body of response from google api")
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return Owner{}, errors.New(string(bodyBytes))
	}

	// Unmarshal response body
	err = json.Unmarshal(bodyBytes, &o)
	if err != nil {
		return Owner{}, errors.New("failed to read body of response from google api")
	}

	return o, nil
}
