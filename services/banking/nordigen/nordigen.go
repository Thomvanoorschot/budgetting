package nordigen

import (
	"budgetting/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	NordigenSecretId    string
	NordigenSecretKey   string
	NordigenUrl         string
	NordigenRedirectUrl string

	accessToken           string
	accessTokenExpiresAt  time.Time
	refreshToken          string
	refreshTokenExpiresAt time.Time
	http.Client
}

func NewClient(config *config.Config) *Client {
	return &Client{
		NordigenRedirectUrl: config.NordigenRedirectUrl,
		NordigenSecretId:    config.NordigenSecretId,
		NordigenSecretKey:   config.NordigenSecretKey,
		NordigenUrl:         config.NordigenUrl,
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
		SecretId:  c.NordigenSecretId,
		SecretKey: c.NordigenSecretKey,
	})

	resp, err := c.Post(fmt.Sprintf("%s/token/new/", c.NordigenUrl), "application/json", bytes.NewBuffer(req))
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
	c.accessToken = tokenResponse.Access
	c.accessTokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.AccessExpires) * time.Second)
	c.refreshToken = tokenResponse.Refresh
	c.refreshTokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.RefreshExpires) * time.Second)
	return nil
}

func (c *Client) refreshAccessToken() error {
	req, _ := json.Marshal(&RefreshAccessTokenRequest{
		Refresh: c.refreshToken,
	})

	resp, err := c.Post(fmt.Sprintf("%s/token/refresh/", c.NordigenUrl), "application/json", bytes.NewBuffer(req))
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
	tokenResponse := &RefreshAccessTokenResponse{}
	err = json.Unmarshal(body, tokenResponse)
	if err != nil {
		return err
	}
	c.accessToken = tokenResponse.Access
	return nil
}

func (c *Client) validateAccessToken() error {
	if c.accessToken == "" {
		return c.getAccessToken()
	}
	if c.accessTokenExpiresAt.Add(time.Duration(10) * time.Minute).Before(time.Now()) {
		return nil
	}
	if c.refreshTokenExpiresAt.Add(time.Duration(2) * time.Minute).Before(time.Now()) {
		return c.refreshAccessToken()
	}
	return c.getAccessToken()
}

type AccessTokenRequest struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}
type RefreshAccessTokenRequest struct {
	Refresh string `json:"refresh"`
}
type AccessTokenResponse struct {
	Access         string `json:"access"`
	AccessExpires  int64  `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int64  `json:"refresh_expires"`
}
type AmountCurrencyPair struct {
	Amount   big.Rat `json:"amount"`
	Currency string  `json:"currency"`
}

type RefreshAccessTokenResponse struct {
	Access        string `json:"access"`
	AccessExpires int64  `json:"access_expires"`
}

type Time time.Time

func (j *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = Time(t)
	return nil
}

func (j *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*j))
}
