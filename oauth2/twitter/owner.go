package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Owner struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
}

type ResponseMe struct {
	Data Owner `json:"data"`
}

type ResponseVerifyCredential struct {
	Email string `json:"email"`
}

func (t *Application) GetOwner(token string) (Owner, error) {
	// GET twitter api
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return Owner{}, errors.New("failed to get owner informations")
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return Owner{}, err
	}
	defer res.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Owner{}, errors.New("failed to read body of response from twitter api")
	}

	// Check status code
	if res.StatusCode != 200 {
		return Owner{}, errors.New(string(bodyBytes))
	}

	// Unmarshal response body
	var resBody ResponseMe
	err = json.Unmarshal(bodyBytes, &resBody)
	if err != nil {
		return Owner{}, errors.New("failed to read body of response from twitter api")
	}

	return resBody.Data, nil
}

func (t *Application) GetOwnerEmail(token string) (string, error) {
	// GET twitter api
	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/account/verify_credentials.json", nil)
	if err != nil {
		return "", errors.New("failed to get owner informations")
	}
	params := req.URL.Query()
	params.Add("include_email", "true")
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("failed to read body of response from twitter api")
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return "", errors.New("failed to get owner informations")
	}

	// Unmarshal response body
	var resBody ResponseVerifyCredential
	err = json.Unmarshal(bodyBytes, &resBody)
	if err != nil {
		return "", errors.New("failed to read body of response from twitter api")
	}

	return resBody.Email, nil
}
