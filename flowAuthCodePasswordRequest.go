package main

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/password/request", handlerPassRequestPost).Methods("POST")
}

func handlerPassRequestPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	code := generateRandomNumbers(9)

	owner := Owner{
		Username:          r.FormValue("username"),
		VerificationCode:  code,
		VerificationRoute: "/ui/password/reset/" + code,
	}

	if id, err := owner.requestPassword(); id != "" && err == nil {
		log.Printf("password reset code for user '%s' is: %s", id, code)
		PassResetStorage.Set(code, id, cache.DefaultExpiration)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusNotFound, ErrorResponse{err.Error()})
	}

	return

}
