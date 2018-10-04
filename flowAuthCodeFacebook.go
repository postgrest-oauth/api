package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"log"
	"net/http"
	"time"
)

var facebookConfig struct {
	ClientId     int64  `env:"OAUTH_FACEBOOK_CLIENT_ID"`
	ClientSecret string `env:"OAUTH_FACEBOOK_CLIENT_SECRET"`
}

func init() {
	err := env.Parse(&facebookConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	Router.HandleFunc("/facebook", handlerFacebookPost).Methods("POST")
}

var fbClient = &http.Client{Timeout: 5 * time.Second}

type FacebookErrorResponse struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FbTraceId string `json:"fbtrace_id"`
}

type FacebookTokenResponse struct {
	AccessToken string                `json:"access_token"`
	Error       FacebookErrorResponse `json:"error"`
}

type FacebookDataResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Hometown  string `json:"hometown"`
	Birthday  string `json:"birthday"`
	Error     FacebookErrorResponse `json:"error"`
}

func getJson(url string, auth_token string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if auth_token != "" {
		req.Header.Set("Authorization", "Bearer "+auth_token)
	}

	if err != nil {
		return err
	}

	r, err := fbClient.Do(req)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	body := r.Body
	return json.NewDecoder(body).Decode(target)
}

func handlerFacebookPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)

	code := r.FormValue("code")
	redirect_uri := r.FormValue("redirect_uri")
	phone := r.FormValue("phone")
	lang := r.FormValue("lang")

	fbTokenUrl := fmt.Sprintf(
		"https://graph.facebook.com/v3.1/oauth/access_token?client_id=%d&redirect_uri=%s&client_secret=%s&code=%s",
		facebookConfig.ClientId,
		redirect_uri,
		facebookConfig.ClientSecret,
		code,
	)
	fbTokenResponse := &FacebookTokenResponse{}
	err := getJson(fbTokenUrl, "", fbTokenResponse)

	if err != nil {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusServiceUnavailable, ErrorResponse{err.Error()})
		return
	}

	errCode := fbTokenResponse.Error.Code
	if errCode > 0 {
		err := fmt.Errorf("facebook token request failed")
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
		return
	}

	fbDataUrl := "https://graph.facebook.com/v3.1/me/?fields=email,first_name,last_name,gender,hometown,birthday"
	fbDataResponse := &FacebookDataResponse{}
	err = getJson(fbDataUrl, fbTokenResponse.AccessToken, fbDataResponse)

	if err != nil {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusServiceUnavailable, ErrorResponse{err.Error()})
		return
	}

	errCode = fbDataResponse.Error.Code
	if errCode > 0 {
		err := fmt.Errorf("facebook data request failed")
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
		return
	}

	json, _ := json.Marshal(fbDataResponse)
	jsonString := string(json)
	owner := Owner{
		FacebookId:   fbDataResponse.Id,
		FacebookJson: jsonString,
	}

	if id, role, jti, err := owner.check_facebook(); err == nil {
		SetSession(id, role, jti, w)
		w.WriteHeader(http.StatusOK)
		return
	}

	if id, role, jti, err := owner.create_or_update_facebook(jsonString, phone, lang); err == nil {
		SetSession(id, role, jti, w)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusConflict, ErrorResponse{err.Error()})
	}


	return
}
