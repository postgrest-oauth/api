package main

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type ClientCredentialsData struct {
	ClientId   string
	ClientRole string
	ClientType string
}

func init() {
	Router.HandleFunc("/token", handlerClientCredentialsToken).
		Methods("POST").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			grantType := r.FormValue("grant_type")
			clientId := r.FormValue("client_id")
			clientSecret := r.FormValue("client_secret")
			if grantType == "client_credentials" && clientId != "" && clientSecret != "" {
				return true
			} else {
				return false
			}
		})
}

func handlerClientCredentialsToken(w http.ResponseWriter, r *http.Request) {
	clientId := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	c := Client{
		Id:     clientId,
		Secret: clientSecret,
	}

	err, ctype := c.check_secret()

	if err != nil {
		e := &errorResponse{Error: "invalid_grant"}
		js, _ := json.Marshal(e)
		jsonResponse(js, w, http.StatusBadRequest)
		return
	} else if ctype != "confidential" {
		e := &errorResponse{Error: "unauthorized_client"}
		js, _ := json.Marshal(e)
		jsonResponse(js, w, http.StatusBadRequest)
		return
	} else {
		data := ClientCredentialsData{
			ClientRole: "msrv-" + clientId,
			ClientId:   clientId,
			ClientType: ctype,
		}
		response := fillClientCredentialsResponse(data)
		js, _ := json.Marshal(response)
		jsonResponse(js, w, http.StatusOK)
		return
	}
}

func fillClientCredentialsResponse(data ClientCredentialsData) tokensResponse {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["type"] = "access_token"
	claims["role"] = data.ClientRole
	claims["client_id"] = data.ClientId
	claims["client_type"] = data.ClientType
	accessTokenString, _ := accessToken.SignedString([]byte(flowConfig.AccessTokenSecret))

	response := tokensResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
	}
	return response
}
