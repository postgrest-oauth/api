package main

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/verify", handlerVerifyPost).Methods("POST")
	Router.HandleFunc("/re-verify", handlerReVerifyPost).Methods("POST")
}

func handlerVerifyPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)

	code := r.FormValue("code")

	savedId, ok := VerifyStorage.Get(code)
	owner := &Owner{}

	if ok {
		owner.Id = savedId.(string)
		if err := owner.verify(); ok && err == nil {
			VerifyStorage.Delete(code)
			log.Printf("user '%s' successfully verified", owner.Id)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Print(err)
			Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
		}
	} else {
		log.Printf("code '%s' is invalid", code)
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{"invalid code"})
	}

	return
}

func handlerReVerifyPost(w http.ResponseWriter, r *http.Request) {
	code := generateRandomNumbers(9)

	owner := Owner{
		Username:          r.FormValue("username"),
		VerificationCode:  code,
		VerificationRoute: authCodeConfig.OauthCodeUi + "/verify/" + code,
	}

	if id, err := owner.reVerify(); err == nil {
		VerifyStorage.Set(code, id, cache.DefaultExpiration)
		log.Printf("Re-verification code for user '%s' is: %s", id, code)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusOK)
	}
	return
}
