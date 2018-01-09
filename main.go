package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"time"
)

type Page struct {
	Owner
	Message          string
	Query            template.URL
	VerificationCode string
}

var ValidateRedirectURI = flag.Bool("validateRedirectURI", true, "Whether validate redirect URI or not. Handy for development")

var tmpl *template.Template
var templatePath = "./templates/"
var Router = mux.NewRouter().StrictSlash(true)
var VerifyStorage = cache.New(24*time.Hour, 2*time.Hour)
var PassResetStorage = cache.New(10*time.Minute, 5*time.Minute)

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client_id")
	redirectUriRequest := r.URL.Query().Get("redirect_uri")
	c := &Client{Id: clientId}
	err, redirectUri := c.check()

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if *ValidateRedirectURI == true {
		if len(redirectUriRequest) > 0 && redirectUri != redirectUriRequest {
			err = errors.New("access denied")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if len(redirectUriRequest) > 0 {
			redirectUri = redirectUriRequest
		}
	}

	ClearSession(w)

	http.Redirect(w, r, redirectUri, 302)
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, templatePath+"favicon.ico")
}

func main() {
	log.Println("Started!")
	flag.Parse()

	blockKeyLength := len(*cookieBlockKey)

	if blockKeyLength != 16 && blockKeyLength != 24 && blockKeyLength != 32 {
		log.Panic("OAUTH_COOKIE_BLOCK_KEY length should be 16, 24 or 32!")
	}

	tmpl = template.Must(template.ParseFiles(
		templatePath+signinTemplate,
		templatePath+signupTemplate,
		templatePath+verifyTemplate,
		templatePath+passwordRequestTemplate,
		templatePath+passwordResetTemplate,
	))

	Router.HandleFunc("/logout", handlerLogout).Methods("GET")
	Router.HandleFunc("/favicon.ico", handlerFavicon)

	corsRouter := cors.AllowAll().Handler(Router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter)
	log.Fatal(http.ListenAndServe(":3684", loggedRouter))
}
