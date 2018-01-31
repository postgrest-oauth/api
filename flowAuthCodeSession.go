package main

import (
	"github.com/caarlos0/env"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
)

var sessionConfig struct {
	CookieHashKey  string `env:"OAUTH_COOKIE_HASH_KEY" envDefault:"supersecret"`
	CookieBlockKey string `env:"OAUTH_COOKIE_BLOCK_KEY" envDefault:"16charssecret!!!"`
}
var cookieHandler *securecookie.SecureCookie

func init() {

	err := env.Parse(&sessionConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	cookieHandler = securecookie.New([]byte(sessionConfig.CookieHashKey), []byte(sessionConfig.CookieBlockKey))

	blockKeyLength := len(sessionConfig.CookieBlockKey)
	if blockKeyLength != 16 && blockKeyLength != 24 && blockKeyLength != 32 {
		log.Panic("COOKIE_BLOCK_KEY length should be 16, 24 or 32!")
	}
}

func SetSession(id string, role string, jti string, response http.ResponseWriter) {
	value := map[string]string{"id": id, "role": role, "jti": jti}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	} else {
		log.Print("Session cookie error")
	}
}

func ClearSession(writer http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(writer, cookie)
}

func GetUser(request *http.Request) (userId string, userRole string, userJti string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userId = cookieValue["id"]
			userRole = cookieValue["role"]
			userJti = cookieValue["jti"]
		}
	}
	return userId, userRole, userJti
}
