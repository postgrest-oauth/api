package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var passwordResetTemplate = "password-reset.html"

func init() {
	Router.HandleFunc("/password/reset/{code}", handlerPassResetGet).Methods("GET").
		Name("verify-pass")
	Router.HandleFunc("/password/reset", handlerPassResetGet).Methods("GET").
		Name("verify-pass-no-code")
	Router.HandleFunc("/password/reset", handlerPassResetPost).Methods("POST")
}

func handlerPassResetGet(w http.ResponseWriter, r *http.Request) {
	s := r.URL.RawQuery
	vars := mux.Vars(r)
	code := vars["code"]
	data := &Page{
		Query:            template.URL(s),
		VerificationCode: code,
		Message:          "WAITING_FOR_CODE",
	}

	err := tmpl.ExecuteTemplate(w, passwordResetTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func handlerPassResetPost(w http.ResponseWriter, r *http.Request) {
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery

	code := r.FormValue("code")
	password := r.FormValue("password")
	data := &Page{
		Query: template.URL(rawQuery),
	}
	savedId, ok := PassResetStorage.Get(code)
	owner := &Owner{}

	if ok {
		owner.Id = savedId.(string)
		owner.Password = password
	}

	if err := owner.resetPassword(); ok && err == nil {
		PassResetStorage.Delete(code)
		data.Message = "VERIFY_SUCCESS"
	} else {
		data.Message = "VERIFY_FAIL"
	}

	err := tmpl.ExecuteTemplate(w, passwordResetTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}
