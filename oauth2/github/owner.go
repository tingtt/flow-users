package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Owner struct {
	Name      string `json:"login"`
	Id        uint64 `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

func (g *Application) GetOwner(token string) (o Owner, err error) {
	// GET github api
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return Owner{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return Owner{}, err
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return Owner{}, errors.New("failed to get owner informations")
	}

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Owner{}, err
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return Owner{}, errors.New(string(bodyBytes))
	}

	// Unmarshal response body
	err = json.Unmarshal(bodyBytes, &o)
	if err != nil {
		return Owner{}, err
	}

	return o, nil
}

type OwnerEmail struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visivility bool   `json:"visibility"`
}

func (g *Application) getOwnerEmails(token string) (emails []OwnerEmail, err error) {
	// GET github api
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bodyBytes))
	}

	// Unmarshal response body
	err = json.Unmarshal(bodyBytes, &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

func (g *Application) GetOwnerPrimaryEmail(token string) (e OwnerEmail, err error) {
	// Get owner emails
	emails, err := g.getOwnerEmails(token)
	if err != nil {
		return OwnerEmail{}, err
	}

	// Find primary email
	found := false
	for _, email := range emails {
		if email.Primary {
			e = email
			found = true
			break
		}
	}
	if !found {
		return OwnerEmail{}, errors.New("primary email not found in body of response from \"api.github.com\"")
	}

	return e, nil
}
