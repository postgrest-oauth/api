package main

import (
	"flag"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
)

var cookieHashKey = flag.String("cookieHashKey", "supersecret", "Hash key for cookie creation. 64 random symbols recommended")
var cookieBlockKey = flag.String("cookieBlockKey", "16charssecret!!!", "Block key for cookie creation. 16, 24 or 32 random symbols are valid")
var cookieHandler = securecookie.New([]byte(*cookieHashKey), []byte(*cookieBlockKey))

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
