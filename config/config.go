package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

const (
	DEVELOPMENT string = "DEVELOPMENT"
	STAGING     string = "STAGING"
	ACCEPTANCE  string = "ACCEPTANCE"
	PRODUCTION  string = "PRODUCTION"
)

type Config struct {
	ApiHost string `envconfig:"API_HOST" default:"localhost"`
	ApiPort string `envconfig:"API_PORT" default:"8080"`

	NordigenRedirectUrl string `envconfig:"NORDIGEN_REDIRECT_URL" default:"https://google.com"`
	NordigenSecretId    string `envconfig:"NORDIGEN_SECRET_ID"    default:"bd4467b1-4beb-4dc5-9e26-5c56cf52b0b1"`
	NordigenSecretKey   string `envconfig:"NORDIGEN_SECRET"       default:"c3de73a380d0bd323f60c499adfce68fe2b8245338d312e00ad0ccdade0b4c920679fedfef6941fed5947a7d03323b73b353a384cee0f70fc7bc07179824eaf6"`
	NordigenUrl         string `envconfig:"NORDIGEN_URL"          default:"https://ob.nordigen.com/api/v2"`

	Auth0ClientId     string `envconfig:"AUTH0_CLIENT_ID"     default:"AzDeO6mWEJrKCFOmOLKYVbBahlh9mKte"`
	Auth0ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET" default:"dkbQ0ocHu4H3Ez5ZRgqx0JdrZESnbSuagUhW-dXB9Az5picHH6NGvM4oyJBJ6CSY"`
	Auth0IssuerUrl    string `envconfig:"AUTH0_ISSUER_URL"    default:"https://dev-msfglgrz.us.auth0.com"`
	Auth0Audience     string `envconfig:"AUTH0_AUDIENCE"      default:"https://budgetting-api"`
	PlanetscaleDSN    string `envconfig:"PlanetscaleDSN"      default:"xj7hji33lvpunmm97fbr:pscale_pw_3jvfPXEImAGvEUpLRrPLUdLUSEOiMzokWUppgaFDcKL@tcp(aws.connect.psdb.cloud)/budgetting?tls=true"`
}

var config *Config
var once sync.Once

// Load reads config file and ENV variables if set.
func Load() *Config {
	once.Do(func() {
		load()
	})

	return config
}

func load() {
	config = new(Config)
	if err := envconfig.Process("", config); err != nil {
		fmt.Println(err)
	}
}
