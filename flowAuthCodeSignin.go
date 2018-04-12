package main

import (
	"log"
	"net/http"
)

func init() {
	Router.HandleFunc("/signin", handlerSigninPost).Methods("POST")
}

func handlerSigninPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)

	owner := Owner{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if id, role, jti, err := owner.check(); err == nil {
		SetSession(id, role, jti, w)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf(err.Error())
		Rnd.JSON(w, http.StatusForbidden, ErrorResponse{err.Error()})
	}

	return
}
