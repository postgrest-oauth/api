package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var verifyTemplate = "verify.html"

func init() {
	Router.HandleFunc("/verify/{code}", handlerVerifyGet).Methods("GET").Name("verify")
	Router.HandleFunc("/verify", handlerVerifyGet).Methods("GET").Name("verify-no-code")
	Router.HandleFunc("/verify", handlerVerifyPost).Methods("POST")
}

func handlerVerifyGet(w http.ResponseWriter, r *http.Request) {
	s := r.URL.RawQuery
	vars := mux.Vars(r)
	code := vars["code"]
	data := &Page{
		Query:            template.URL(s),
		VerificationCode: code,
		Message:          "WAITING_FOR_CODE",
	}

	err := tmpl.ExecuteTemplate(w, verifyTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func handlerVerifyPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery

	code := r.FormValue("code")
	data := &Page{
		Query: template.URL(rawQuery),
	}
	savedId, ok := VerifyStorage.Get(code)
	owner := &Owner{}

	if ok {
		owner.Id = savedId.(string)
	}

	if err := owner.verify(); ok && err == nil {
		VerifyStorage.Delete(code)
		data.Message = "VERIFY_SUCCESS"
	} else {
		data.Message = "VERIFY_FAIL"
	}

	err := tmpl.ExecuteTemplate(w, verifyTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}
