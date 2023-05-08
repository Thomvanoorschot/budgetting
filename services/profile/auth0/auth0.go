package auth0

import (
	"budgetting/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Auth0ClientId     string
	Auth0ClientSecret string
	Auth0IssuerUrl    string

	accessToken          string
	accessTokenExpiresAt time.Time
	http.Client
}

func NewClient(config *config.Config) *Client {
	return &Client{
		Auth0ClientId:     config.Auth0ClientId,
		Auth0ClientSecret: config.Auth0ClientSecret,
		Auth0IssuerUrl:    config.Auth0IssuerUrl,
		Client: http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
	}
}

func (c *Client) Execute(req *http.Request, o interface{}) error {
	err := c.validateAccessToken()
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", c.accessToken)},
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, o)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getAccessToken() error {
	req, _ := json.Marshal(&AccessTokenRequest{
		ClientId:     c.Auth0ClientId,
		ClientSecret: c.Auth0ClientSecret,
		Audience:     fmt.Sprintf("%s/api/v2/", c.Auth0IssuerUrl),
		GrantType:    "client_credentials",
	})

	resp, err := c.Post(fmt.Sprintf("%s/oauth/token", c.Auth0IssuerUrl), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	tokenResponse := &AccessTokenResponse{}
	err = json.Unmarshal(body, tokenResponse)
	if err != nil {
		return err
	}
	c.accessToken = tokenResponse.AccessToken
	c.accessTokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	return nil
}

func (c *Client) validateAccessToken() error {
	if c.accessToken == "" || c.accessTokenExpiresAt.Add(time.Duration(10)*time.Minute).After(time.Now()) {
		return c.getAccessToken()
	}
	return nil
}

type AccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
