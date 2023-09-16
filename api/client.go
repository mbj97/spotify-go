package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Status int

const (
	STATUS_WAITING_FOR_AUTH Status = iota
	STATUS_AUTHORIZED
)

type SpotifyClient struct {
	Client   *spotify.Client
	username string
	LoginURL string
	Status   Status
}

const (
	redirectURI = "http://localhost:8080"
	tokenPath   = "tokens/%s.json"
)

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadCurrentlyPlaying),
	)
)

func CreateNewClient(username string) (*SpotifyClient, error) {
	token, err := loadTokenFromFile(username)
	spotifyClient := SpotifyClient{username: username}
	if err == nil {
		// Use the saved access token if it exists and is not expired
		fmt.Println("TOKEN VALID")
		if token.Valid() {
			client := spotify.New(auth.Client(context.Background(), token))
			spotifyClient.Client = client
			spotifyClient.Status = STATUS_AUTHORIZED
			return &spotifyClient, nil
		}
	}
	// No saved/expired token, create login link
	url := auth.AuthURL(username)
	spotifyClient.LoginURL = url
	spotifyClient.Status = STATUS_WAITING_FOR_AUTH
	fmt.Println(url + "@@")

	return &spotifyClient, nil
}

// func NewSpotifyClient(state string) (*SpotifyClient, error) {
// 	// Check if a token file exists
// 	token, err := loadTokenFromFile(tokenFile)
// 	if err == nil {
// 		// Use the saved access token if it exists and is not expired
// 		if token.Valid() {
// 			client := spotify.New(auth.Client(context.Background(), token))
// 			spotifyClient := SpotifyClient{Client: client}
// 			return &spotifyClient, nil
// 		}
// 	}

// 	// No valid token found, proceed with the authorization flow
// 	http.HandleFunc("/", completeAuth)
// 	go func() {
// 		err := http.ListenAndServe(":8080", nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	url := auth.AuthURL(state)
// 	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

// 	client := <-ch
// 	cToken, err := client.Token()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	saveTokenToFile(tokenFile, cToken)

// 	user, err := client.CurrentUser(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("You are logged in as:", user.ID)

// 	spotifyClient := SpotifyClient{Client: client}
// 	return &spotifyClient, nil
// }

func (sc *SpotifyClient) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), sc.username, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		fmt.Fprintf(w, "Error: %s", err.Error())
	}
	if st := r.FormValue("state"); st != sc.username {
		http.NotFound(w, r)
		fmt.Fprintf(w, "State mismatch: %s != %s\n", st, sc.username)
	}
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	// ch <- client
	cToken, err := client.Token()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}
	// Save client
	sc.Client = client
	sc.Status = STATUS_AUTHORIZED
	sc.SaveTokenToFile(cToken)
}

// loadTokenFromFile Checks if token file exists and loads it
func loadTokenFromFile(username string) (*oauth2.Token, error) {
	filePath := fmt.Sprintf(tokenPath, username)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}
	file, err := os.Open(filePath)
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
func (sc SpotifyClient) SaveTokenToFile(token *oauth2.Token) error {
	filePath := fmt.Sprintf(tokenPath, sc.username)
	file, err := os.Create(filePath)
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
