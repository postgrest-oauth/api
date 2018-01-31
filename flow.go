package main

import (
	"log"
	"github.com/caarlos0/env"
	"net/http"
)

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	State            string `json:"state,omitempty"`
}

var flowConfig struct {
	AccessTokenSecret  string `env:"OAUTH_ACCESS_TOKEN_SECRET" envDefault:"morethan32symbolssecretkey!!!!!!"`
	AccessTokenTTL     int    `env:"OAUTH_ACCESS_TOKEN_TTL" envDefault:"7200"`
	RefreshTokenSecret string `env:"OAUTH_REFRESH_TOKEN_SECRET" envDefault:"notlesshan32symbolssecretkey!!!!"`
}

func init() {
	err := env.Parse(&flowConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}

}

func jsonResponse(js []byte, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}