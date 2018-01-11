package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

type Data struct {
	ClientId string
	UserId   string
	UserRole string
}

var Storage = cache.New(10*time.Minute, 20*time.Minute)

var AccessTokenSecret = flag.String("accessTokenJWTSecret", "morethan32symbolssecretkey!!!!!!",
	"Secret key for generating JWT access tokens")
var AccessTokenTTL = flag.Int64("accessTokenTTL", 7200, "Access token TTL in seconds")
var RefreshTokenSecret = flag.String("refreshTokenJWTSecret", "notlesshan32symbolssecretkey!!!!",
	"Secret key for generating JWT refresh tokens")

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

	if *ValidateRedirectURI == true {
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

	uId, uRole := GetUser(r)
	if uId == "" || uRole == "" {
		http.Redirect(w, r, "/signin?"+r.URL.RawQuery, 302)
		return
	}

	code := generateRandomString(randNumbers, 9)

	data := &Data{ClientId: clientId, UserId: uId, UserRole: uRole}

	Storage.Set(code, *data, cache.DefaultExpiration)

	redirectUri = redirectUri + "?code=" + code
	if state := r.URL.Query().Get("state"); state != "" {
		redirectUri = redirectUri + "&state=" + state
	}

	http.Redirect(w, r, redirectUri, 302)
	return
}

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	State            string `json:"state,omitempty"`
}

func handlerAuthCodeToken(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	clientId := r.FormValue("client_id")

	if data, ok := Storage.Get(code); ok {
		if data.(Data).ClientId == clientId {
			response := fillTokensResponse(data.(Data))
			js, _ := json.Marshal(response)
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
		return []byte(*RefreshTokenSecret), nil
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
		o := Owner{Id: userId}
		errRole, role := o.getOwnerRoleById()
		if errRole != nil {
			e := &errorResponse{Error: "invalid_grant"}
			js, _ := json.Marshal(e)
			jsonResponse(js, w, http.StatusBadRequest)
			return
		}
		d := Data{
			UserId:   userId,
			ClientId: clientId,
			UserRole: role,
		}

		response := fillTokensResponse(d)
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

func fillTokensResponse(data Data) tokensResponse {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["type"] = "access_token"
	claims["role"] = data.UserRole
	claims["id"] = data.UserId
	claims["client_id"] = data.ClientId
	claims["exp"] = time.Now().Add(time.Second * time.Duration(*AccessTokenTTL)).Unix()
	accessTokenString, _ := accessToken.SignedString([]byte(*AccessTokenSecret))

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["type"] = "refresh_token"
	claims["id"] = data.UserId
	claims["client_id"] = data.ClientId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	refreshTokenString, _ := refreshToken.SignedString([]byte(*RefreshTokenSecret))

	response := tokensResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "bearer",
	}
	return response
}

func jsonResponse(js []byte, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}
