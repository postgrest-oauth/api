package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var signinTemplate = "signin.html"

func init() {
	Router.HandleFunc("/signin", handlerSigninGet).Methods("GET").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			if strings.Contains(r.RequestURI, "/signin?response_type=") {
				return true
			} else {
				return false
			}
		})
	Router.HandleFunc("/signin", handlerSigninPost).Methods("POST")
}

func handlerSigninGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Query: template.URL(s[8:]),
	}
	err := tmpl.ExecuteTemplate(w, signinTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerSigninPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery

	data := &Page{
		Owner: Owner{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		},
		Query: template.URL(rawQuery),
	}

	err, id, role := data.Owner.check()

	if err != nil {
		data.Message = "WRONG_CREDENTIALS"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, signinTemplate, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		http.Redirect(w, r, "/authorize?"+rawQuery, 302)
		return
	}
}
