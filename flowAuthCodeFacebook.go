package main

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

var facebookConfig struct {
	ClientId     string `env:"OAUTH_FACEBOOK_CLIENT_ID" envDefault:"supersecret"`
	ClientSecret string `env:"OAUTH_FACEBOOK_CLIENT_SECRET" envDefault:"16charssecret!!!"`
}

func init() {

	err := env.Parse(&facebookConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	Router.HandleFunc("/facebook/url", handlerFacebookUrl).Methods("GET")
	Router.HandleFunc("/facebook/enter", handlerFacebookEnter).Methods("POST")
}

func handlerFacebookUrl(w http.ResponseWriter, r *http.Request) {

	redirectUri := r.URL.Query().Get("redirect_uri")

	if redirectUri == "" {
		redirectUri = authCodeConfig.OauthCodeUi
	}

	authURL, err := gocial.New().
		Driver("facebook").        // Set provider
		Scopes([]string{"email"}). // Set optional scope(s)
		Redirect(                  //
			facebookConfig.ClientId, // Client ID
			facebookConfig.ClientSecret,
			redirectUri, // Redirect URL
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		Rnd.JSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	} else {
		type Response struct {
			Url string `json:"url"`
		}
		js, _ := json.Marshal(Response{authURL})
		jsonResponse(js, w, http.StatusOK)
		return
	}

	return
}
func handlerFacebookEnter(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	state := r.FormValue("state")

	// Handle callback and check for errors
	data, _, err := gocial.Handle(state, code)
	if err != nil {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
		return
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
		return
	}

	owner := Owner{FacebookId: data.ID, Data: string(dataJson)}

	if id, role, jti, err := owner.checkFacebook(); err == nil {

		if id != "" {
			SetSession(id, role, jti, w)
			w.WriteHeader(http.StatusOK)
		} else {
			if id, role, jti, err := owner.createFacebook(); err == nil {
				SetSession(id, role, jti, w)
				w.WriteHeader(http.StatusOK)
			} else {
				log.Printf(err.Error())
				Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
			}
		}

	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
	}

	return
}
