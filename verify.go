package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var verifyTemplate = "verify.html"

func init() {
	Router.HandleFunc("/verify/{id:[0-9]+}/{code}", handlerVerifyGet).Methods("GET").Name("verify")
	Router.HandleFunc("/verify/{id:[0-9]+}", handlerVerifyGet).Methods("GET")
	Router.HandleFunc("/verify/{id:[0-9]+}", handlerVerifyPost).Methods("POST")
}

func handlerVerifyGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	vars := mux.Vars(r)
	id := vars["id"]
	code := vars["code"]
	data := &Page{
		Query: template.URL(s[1:]),
	}
	owner := &Owner{Id: id}
	err := owner.verify()

	if code == "" {
		data.Message = "WAITING_FOR_CODE"
	} else if savedId, ok := VerifyStorage.Get(code); ok && id == savedId && err == nil {
		VerifyStorage.Delete(code)
		data.Message = "VERIFY_SUCCESS"
	} else {
		data.Message = "VERIFY_FAIL"
	}

	err = tmpl.ExecuteTemplate(w, verifyTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func handlerVerifyPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	code := r.FormValue("code")
	data := &Page{}
	owner := &Owner{Id: id}
	err := owner.verify()

	if savedId, ok := VerifyStorage.Get(code); ok && id == savedId && err == nil {
		VerifyStorage.Delete(code)
		data.Message = "VERIFY_SUCCESS"
	} else {
		data.Message = "VERIFY_FAIL"
	}

	err = tmpl.ExecuteTemplate(w, verifyTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}
