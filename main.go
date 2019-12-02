package main

import (
	"github.com/caarlos0/env"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var Router = mux.NewRouter().StrictSlash(true)

var authCodeConfig struct {
	ValidateRedirectURI bool   `env:"OAUTH_VALIDATE_REDIRECT_URI" envDefault:"true"`
	OauthCodeUi         string `env:"OAUTH_CODE_UI" envDefault:"http://localhost:3685"`
}

func init() {
	err := env.Parse(&authCodeConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}
}

func main() {
	log.Println("Started!")
	corsRouter := cors.New(cors.Options{
		AllowedOrigins: []string{"http://foo.com"},
	}).Handler(Router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter)
	log.Fatal(http.ListenAndServe(":3684", loggedRouter))
}
