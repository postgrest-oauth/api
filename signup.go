package main

import (
	"github.com/patrickmn/go-cache"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var signupTemplate = "signup.html"

func init() {
	Router.HandleFunc("/signup", handlerSignupGet).Methods("GET")
	Router.HandleFunc("/signup", handlerSignupPost).Methods("POST")
}

func handlerSignupGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Query: template.URL(s[8:]),
	}
	err := tmpl.ExecuteTemplate(w, signupTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerSignupPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	s := r.RequestURI
	code := generateRandomString(randNumbers, 6)

	data := &Page{
		Owner: Owner{
			Email:            r.FormValue("email"),
			Phone:            r.FormValue("phone"),
			Password:         r.FormValue("password"),
			VerificationCode: code,
		},
		Query: template.URL(s[8:]),
	}
	err, id, role := data.Owner.create()

	if err != nil {
		data.Message = "WRONG_INPUT"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, signupTemplate, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		route, _ := Router.Get("verify").URL("id", id, "code", code)
		routeNoCode, _ := Router.Get("verify-no-code").URL("id", id)
		data.Owner.VerificationRoute = route.String()
		log.Printf("Verification route for user '%s' is: %s", id, route.String())
		VerifyStorage.Set(code, id, cache.DefaultExpiration)
		refUrl, _ := url.Parse(r.Referer())
		rawQuery := refUrl.RawQuery
		http.Redirect(w, r, routeNoCode.String()+"?"+rawQuery, 302)
		return
	}

}
