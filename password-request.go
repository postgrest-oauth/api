package main

import (
	"github.com/patrickmn/go-cache"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var passwordRequestTemplate = "password-request.html"

func init() {
	Router.HandleFunc("/password/request", handlerPassRequestGet).Methods("GET")
	Router.HandleFunc("/password/request", handlerPassRequestPost).Methods("POST")
}

func handlerPassRequestGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Query: template.URL(s[1:]),
	}
	err := tmpl.ExecuteTemplate(w, passwordRequestTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerPassRequestPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	code := generateRandomNumbers(9)

	data := &Page{
		Owner: Owner{
			Username:         r.FormValue("username"),
			VerificationCode: code,
		},
	}
	_, id := data.Owner.requestPassword()

	route, _ := Router.Get("verify-pass").URL("code", code)
	routeNoCode, _ := Router.Get("verify-pass-no-code").URL()
	data.Owner.VerificationRoute = route.String()
	if id != "" {
		log.Printf("password reset route for user '%s' is: %s", id, route.String())
		PassResetStorage.Set(code, id, cache.DefaultExpiration)
	} else {
		log.Printf("password reset for user '%s' failed. User not found", data.Owner.Username)
	}
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery
	http.Redirect(w, r, routeNoCode.String()+"?"+rawQuery, 302)

	return

}
