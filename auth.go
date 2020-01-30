package main

import (
	"fmt"
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
	http.HandleFunc("/callback", completeAuth)
	// start our server
	go http.ListenAndServe(":8082", nil)

	// create the auth url based on our env vars
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for the auth channel to complete
	client := <-ch
	// use the client to make calls that require authorization
	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
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
