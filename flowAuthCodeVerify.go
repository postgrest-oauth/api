package main

import (
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/verify", handlerVerifyPost).Methods("POST")
}

func handlerVerifyPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)

	code := r.FormValue("code")

	savedId, ok := VerifyStorage.Get(code)
	owner := &Owner{}

	if ok {
		owner.Id = savedId.(string)
	}

	if err := owner.verify(); ok && err == nil {
		VerifyStorage.Delete(code)
		log.Printf("user '%s' successfully verified", owner.Id)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Print(err)
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
	}

	return
}
