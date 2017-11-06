package main

import (
	"net/http"
	"log"
)

func SetSession(id string, role string, response http.ResponseWriter) {
	value := map[string]string{"id": id, "role": role}
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

func GetUser(request *http.Request) (userId string, userRole string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userId = cookieValue["id"]
			userRole = cookieValue["role"]
		}
	}
	return userId, userRole
}
