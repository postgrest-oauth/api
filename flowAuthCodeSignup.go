package main

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/signup", handlerSignupPost).Methods("POST")
}

func handlerSignupPost(w http.ResponseWriter, r *http.Request) {
	code := generateRandomNumbers(9)

	owner := Owner{
		Email:             r.FormValue("email"),
		Phone:             r.FormValue("phone"),
		Password:          r.FormValue("password"),
		Language:          r.FormValue("language"),
		VerificationCode:  code,
		VerificationRoute: authCodeConfig.OauthCodeUi + "/verify/" + code,
	}

	if id, err := owner.create(); err == nil {
		VerifyStorage.Set(code, id, cache.DefaultExpiration)
		log.Printf("Verification code for user '%s' is: %s", id, code)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
	}
	return
}
