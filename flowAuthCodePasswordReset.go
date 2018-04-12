package main

import (
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/ui/password/reset", handlerPassResetPost).Methods("POST")
}

func handlerPassResetPost(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	password := r.FormValue("password")

	savedId, ok := PassResetStorage.Get(code)
	owner := &Owner{}

	if ok {
		owner.Id = savedId.(string)
		owner.Password = password
	}

	if err := owner.resetPassword(); err == nil {
		PassResetStorage.Delete(code)
		log.Printf("password reset for user '%s' was successful", owner.Id)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
	}

	return
}
