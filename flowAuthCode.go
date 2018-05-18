package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"net/url"
)

type AuthCodeData struct {
	ClientId   string
	ClientType string
	UserId     string
	UserRole   string
	UserJti    string
}

var Storage = cache.New(10*time.Minute, 20*time.Minute)

func init() {
	Router.HandleFunc("/authorize", handlerAuthCode).
		Methods("GET").
		Queries("response_type", "code", "client_id", "{client_id}")

	Router.HandleFunc("/token", handlerAuthCodeToken).
		Methods("POST").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			grantType := r.FormValue("grant_type")
			code := r.FormValue("code")
			clientId := r.FormValue("client_id")
			if grantType == "authorization_code" && code != "" && clientId != "" {
				return true
			} else {
				return false
			}
		})

	Router.HandleFunc("/token", handlerAuthCodeRefreshToken).
		Methods("POST").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			grantType := r.FormValue("grant_type")
			refreshToken := r.FormValue("refresh_token")
			if grantType == "refresh_token" && refreshToken != "" {
				return true
			} else {
				return false
			}
		})
}

func handlerAuthCode(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client_id")
	redirectUriRequest := r.URL.Query().Get("redirect_uri")
	c := &Client{Id: clientId}
	err, redirectUri := c.check()

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if authCodeConfig.ValidateRedirectURI == true {
		if len(redirectUriRequest) > 0 && redirectUri != redirectUriRequest {
			err = errors.New("access denied")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if len(redirectUriRequest) > 0 {
			redirectUri = redirectUriRequest
		}
	}

	uId, uRole, uJti := GetUser(r)
	if uId == "" || uRole == "" {
		http.Redirect(w, r, authCodeConfig.OauthCodeUi + "/signin?"+r.URL.RawQuery, 302)
		return
	}

	code := generateRandomString(24)

	data := &AuthCodeData{ClientId: clientId, UserId: uId, UserRole: uRole, UserJti: uJti, ClientType: "public"}

	Storage.Set(code, *data, cache.DefaultExpiration)

	redirectUriParsed, _ := url.Parse(redirectUri)
	params, _ := url.ParseQuery(redirectUriParsed.RawQuery)
	params.Add("code", code)

	if state := r.URL.Query().Get("state"); state != "" {
		params.Add("state", state)
	}

	http.Redirect(w, r, redirectUriParsed.String(), 302)
	return
}

func handlerAuthCodeToken(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	clientId := r.FormValue("client_id")

	if data, ok := Storage.Get(code); ok {
		if data.(AuthCodeData).ClientId == clientId {
			response := fillAuthFlowResponse(data.(AuthCodeData))
			js, _ := json.Marshal(response)
			Storage.Delete(code)
			jsonResponse(js, w, http.StatusOK)
			return
		} else {
			err := errors.New("wrong client id")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		e := &errorResponse{Error: "invalid_grant"}
		js, _ := json.Marshal(e)

		jsonResponse(js, w, http.StatusBadRequest)
		return
	}
}

func handlerAuthCodeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenString := r.FormValue("refresh_token")
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(flowConfig.RefreshTokenSecret), nil
	})

	if err != nil {
		log.Print(err)
		e := &errorResponse{Error: "invalid_grant", ErrorDescription: err.Error()}
		js, _ := json.Marshal(e)

		jsonResponse(js, w, http.StatusBadRequest)
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid && claims["type"] == "refresh_token" {
		userId := claims["id"].(string)
		clientId := claims["client_id"].(string)
		claimsJti := claims["jti"].(string)
		claimsRole := claims["role"].(string)
		o := Owner{Id: userId}
		role, jti, err := o.getOwnerRoleAndJtiById()
		if err != nil {
			e := &errorResponse{Error: "invalid_grant"}
			js, _ := json.Marshal(e)
			jsonResponse(js, w, http.StatusBadRequest)
			return
		}
		if role != claimsRole {
			e := &errorResponse{Error: "invalid_grant"}
			js, _ := json.Marshal(e)
			jsonResponse(js, w, http.StatusBadRequest)
			return
		}
		if jti != claimsJti {
			e := &errorResponse{Error: "invalid_grant"}
			js, _ := json.Marshal(e)
			jsonResponse(js, w, http.StatusBadRequest)
			return
		}
		d := AuthCodeData{
			UserId:   userId,
			ClientId: clientId,
			UserRole: role,
			UserJti:  jti,
		}

		response := fillAuthFlowResponse(d)
		js, _ := json.Marshal(response)
		jsonResponse(js, w, http.StatusOK)
		return
	} else {
		e := &errorResponse{Error: "invalid_grant", ErrorDescription: "token is invalid"}
		js, _ := json.Marshal(e)
		jsonResponse(js, w, http.StatusBadRequest)
		return
	}
}

func fillAuthFlowResponse(data AuthCodeData) tokensResponse {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["type"] = "access_token"
	claims["role"] = data.UserRole
	claims["id"] = data.UserId
	claims["client_id"] = data.ClientId
	claims["client_type"] = data.ClientType
	claims["jti"] = data.UserJti
	claims["exp"] = time.Now().Add(time.Second * time.Duration(flowConfig.AccessTokenTTL)).Unix()
	accessTokenString, _ := accessToken.SignedString([]byte(flowConfig.AccessTokenSecret))

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["type"] = "refresh_token"
	claims["id"] = data.UserId
	claims["client_id"] = data.ClientId
	claims["client_type"] = data.ClientType
	claims["role"] = data.UserRole
	claims["jti"] = data.UserJti
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	refreshTokenString, _ := refreshToken.SignedString([]byte(flowConfig.RefreshTokenSecret))

	response := tokensResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "bearer",
	}
	return response
}
