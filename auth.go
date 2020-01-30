package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8082/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserTopRead, spotify.ScopePlaylistModifyPublic)
	ch    = make(chan *spotify.Client)
	state = strconv.Itoa(int(time.Now().Unix()))
)

func login() *spotify.Client {
	// First we load our environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	auth.SetAuthInfo(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	// Next we start an HTTP server to listen for authentication
	http.HandleFunc("/callback", handleAuth)
	http.HandleFunc("/", handleIndex)
	// start our server
	go http.ListenAndServe(":8082", nil)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", auth.AuthURL(state))
	// wait for the auth channel to complete
	client := <-ch
	// use the client to make calls that require authorization
	return client
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed! You may close this window")
	ch <- &client
}

// HTMLLink is a struct for creating a link with a dynamic href attribute using html/templating
type HTMLLink struct {
	HREF string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("authentication url")
	tmpl, _ = tmpl.Parse("<a href='{{.HREF}}'>Auth with Spotify</a>")
	l := HTMLLink{HREF: auth.AuthURL(state)}
	tmpl.Execute(w, l)
}
