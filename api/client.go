package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyClient struct {
	C *spotify.Client
}

const (
	redirectURI = "http://localhost:8080"
	state       = "abc123"
	tokenFile   = "spotify_token.json"
)

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadCurrentlyPlaying),
	)
	ch = make(chan *spotify.Client)
)

func NewSpotifyClient() (*SpotifyClient, error) {
	// Check if a token file exists
	token, err := loadTokenFromFile(tokenFile)
	if err == nil {
		// Use the saved access token if it exists and is not expired
		if token.Valid() {
			client := spotify.New(auth.Client(context.Background(), token))
			spotifyClient := SpotifyClient{C: client}
			return &spotifyClient, nil
		}
	}

	// No valid token found, proceed with the authorization flow
	http.HandleFunc("/", completeAuth)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch
	cToken, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}
	saveTokenToFile(tokenFile, cToken)

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("You are logged in as:", user.ID)

	spotifyClient := SpotifyClient{C: client}
	return &spotifyClient, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
	cToken, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}
	saveTokenToFile(tokenFile, cToken)
}

// loadTokenFromFile Checks if token file exists and loads it
func loadTokenFromFile(filename string) (*oauth2.Token, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := new(oauth2.Token)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// saveTokenToFile Creates or overwrites a new token file
func saveTokenToFile(filename string, token *oauth2.Token) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(token); err != nil {
		return err
	}

	return nil
}
