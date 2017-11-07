package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"net/url"
	"flag"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/gorilla/securecookie"
	"github.com/rs/cors"
)

type Page struct {
	Owner
	Title   string
	Message string
	Action  string // Can be "signin" or "signup"
	Query template.URL
}

var tmpl *template.Template
var cookieHandler *securecookie.SecureCookie

var templatePath = flag.String("templatePath", "./", "Path to template html file. With trailing slash")
var templateName = flag.String("templateName", "index.html", "Name of template html file")

var cookieHashKey = flag.String("cookieHashKey", "supersecret", "Hash key for cookie creation. 64 random symbols recommended")
var cookieBlockKey = flag.String("cookieBlockKey", "16charssecret!!!", "Block key for cookie creation. 16, 24 or 32 random symbols are valid")

var dbConnString = flag.String("dbConnString", "postgres://user:pass@localhost:5432/test?sslmode=disable",
	"Database connection string")

var AccessTokenSecret = flag.String("accessTokenJWTSecret", "morethan32symbolssecretkey!!!!!!",
	"Secret key for generating JWT access tokens")
var RefreshTokenSecret = flag.String("refreshTokenJWTSecret", "notlesshan32symbolssecretkey!!!!",
	"Secret key for generating JWT refresh tokens")

	var Router = mux.NewRouter().StrictSlash(true)

func handlerHomeGet(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	data := &Page{
		Title: "SignIn/SignUp",
		Query: template.URL(s[2:]),
	}
	err := tmpl.ExecuteTemplate(w, *templateName, data)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handlerHomePost(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	refUrl, _ := url.Parse(r.Referer())
	rawQuery := refUrl.RawQuery

	data := &Page{
		Owner: Owner{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		},
		Title: "SignIn/SignUp",
		Query: template.URL(rawQuery),
	}

	err, id, role := data.Owner.check()

	if err != nil {
		data.Message = "Wrong username or password"
		data.Action = "signin"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, *templateName, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		http.Redirect(w, r, "/authorize?" + rawQuery, 302)
		return
	}
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
	s := r.RequestURI
	ClearSession(w)
	data := &Page{
		Owner: Owner{
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: r.FormValue("password"),
		},
		Title: "SignIn/SignUp",
		Query: template.URL(s[8:]),
	}
	err, id, role := data.Owner.create()

	if err != nil {
		data.Message = "User can't be created"
		data.Action = "signup"
		data.Owner.Password = ""
		err = tmpl.ExecuteTemplate(w, *templateName, data)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		SetSession(id, role, w)
		refUrl, _ := url.Parse(r.Referer())
		rawQuery := refUrl.RawQuery
		http.Redirect(w, r, "/authorize?" + rawQuery, 302)
		return
	}

}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	http.Redirect(w, r, "/", 302)
}

func handlerFavicon (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, *templatePath + "favicon.ico")
}

func main() {
	log.Println("Started!")
	flag.Parse()

	blockKeyLength := len(*cookieBlockKey)

	if blockKeyLength != 16 && blockKeyLength != 24 && blockKeyLength != 32 {
		log.Panic("OAUTH_COOKIE_BLOCK_KEY length should be 16, 24 or 32!")
	}

	tmpl = template.Must(template.ParseFiles(*templatePath + *templateName))
	cookieHandler = securecookie.New([]byte(*cookieHashKey), []byte(*cookieBlockKey))

	Router.HandleFunc("/", handlerHomeGet).Methods("GET").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		if strings.Contains(r.RequestURI, "/?response_type=") {
			return true
		} else {
			return false
		}
	})

	Router.HandleFunc("/", handlerHomePost).Methods("POST")
	Router.HandleFunc("/signup", handlerSignup).Methods("POST", "GET")
	Router.HandleFunc("/logout", handlerLogout).Methods("GET")

	Router.HandleFunc("/favicon.ico", handlerFavicon)

	corsRouter := cors.AllowAll().Handler(Router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter)
	log.Fatal(http.ListenAndServe(":3684", loggedRouter))
}
