package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

var facebookConfig struct {
	ClientId        string `env:"OAUTH_FB_CLIENT_ID"`
	ClientSecret    string `env:"OAUTH_FB_CLIENT_SECRET"`
	RedirectUriHost string `env:"OAUTH_FB_REDIRECT_URI_HOST" envDefault:"http://localhost:3684"`
	ApiVersion      string `env:"OAUTH_FB_API_VERSION" envDefault:"v2.12"`
}

func init() {
	err := env.Parse(&facebookConfig)
	if err != nil {
		log.Printf("%+v\n", err)

	}
	Router.HandleFunc("/facebook/login", handlerFacebookLoginGet).Methods("GET")
	Router.HandleFunc("/facebook/cb", handlerFacebookLoginGet).Methods("GET").Name("facebook-callback")
}

func handlerFacebookLoginGet(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(5)
	cbRoute, _ := Router.Get("facebook-callback").URL()
	fbLink := fmt.Sprintf(
		"https://www.facebook.com/%s/dialog/oauth?client_id=%s&redirect_uri=%s&state=%s",
		facebookConfig.ApiVersion,
		facebookConfig.ClientId,
		facebookConfig.RedirectUriHost + cbRoute.String(),
		state,
	)
	http.Redirect(w, r, fbLink, 302)
	return
}
