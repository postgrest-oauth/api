package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"github.com/satori/go.uuid"
	"time"
)

type Page struct {
	Owner
	Message     string
	Query       template.URL
}

var tmpl *template.Template
var cookieHandler *securecookie.SecureCookie

var templatePath = "./templates/"
var signinTemplate = "signin.html"
var signupTemplate = "signup.html"
var verifyTemplate = "verify.html"

var cookieHashKey = flag.String("cookieHashKey", "supersecret", "Hash key for cookie creation. 64 random symbols recommended")
var cookieBlockKey = flag.String("cookieBlockKey", "16charssecret!!!", "Block key for cookie creation. 16, 24 or 32 random symbols are valid")

var ValidateRedirectURI = flag.Bool("validateRedirectURI", true, "Whether validate redirect URI or not. Handy for development")

var Router = mux.NewRouter().StrictSlash(true)

var VerifyStorage = cache.New(24*time.Hour, 2*time.Hour)

func handlerSigninGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Query: template.URL(s[8:]),
	}
	err := tmpl.ExecuteTemplate(w, signinTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerSigninPost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery

	data := &Page{
		Owner: Owner{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		},
		Query: template.URL(rawQuery),
	}

	err, id, role := data.Owner.check()

	if err != nil {
		data.Message = "WRONG_CREDENTIALS"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, signinTemplate, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		http.Redirect(w, r, "/authorize?"+rawQuery, 302)
		return
	}
}

func handlerSignupGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Query: template.URL(s[8:]),
	}
	err := tmpl.ExecuteTemplate(w, signupTemplate, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerSignupPost(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	code := uuid.NewV4().String()
	ClearSession(w)
	data := &Page{
		Owner: Owner{
			Email:            r.FormValue("email"),
			Phone:            r.FormValue("phone"),
			Password:         r.FormValue("password"),
			VerificationCode: code,
		},
		Query: template.URL(s[8:]),
	}
	err, id, role := data.Owner.create()

	if err != nil {
		data.Message = "WRONG_INPUT"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, signupTemplate, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		log.Printf("Code for user '%s' is: %s", id, code)
		VerifyStorage.Set(code, id, cache.DefaultExpiration)
		refUrl, _ := url.Parse(r.Referer())
		rawQuery := refUrl.RawQuery
		http.Redirect(w, r, "/authorize?"+rawQuery, 302)
		return
	}

}

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

func handlerVerify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	code := vars["code"]
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

	tmpl = template.Must(template.ParseFiles(templatePath+signinTemplate, templatePath+signupTemplate, templatePath+verifyTemplate))
	cookieHandler = securecookie.New([]byte(*cookieHashKey), []byte(*cookieBlockKey))

	Router.HandleFunc("/signin", handlerSigninGet).Methods("GET").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			if strings.Contains(r.RequestURI, "/signin?response_type=") {
				return true
			} else {
				return false
			}
		})

	Router.HandleFunc("/signin", handlerSigninPost).Methods("POST")
	Router.HandleFunc("/signup", handlerSignupGet).Methods("GET")
	Router.HandleFunc("/signup", handlerSignupPost).Methods("POST")
	Router.HandleFunc("/logout", handlerLogout).Methods("GET")

	Router.HandleFunc("/verify/{id:[0-9]+}/{code}", handlerVerify).Methods("GET")

	Router.HandleFunc("/favicon.ico", handlerFavicon)

	corsRouter := cors.AllowAll().Handler(Router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter)
	log.Fatal(http.ListenAndServe(":3684", loggedRouter))
}
